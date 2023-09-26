use rand::Rng;
use rayon::prelude::*;
use sha2::{Digest, Sha512};
use std::arch::x86_64::_mm_stream_si64;
use std::collections::HashSet;
use std::io::{BufRead, BufReader};
use std::iter::zip;
const NTT_DEG: u32 = 26;
struct Ntt<const MOD: u64, const ROOT: u64>;
impl<const MOD: u64, const ROOT: u64> Ntt<MOD, ROOT> {
    fn add(x: u64, y: u64) -> u64 {
        let s = x + y;
        if s >= MOD {
            s - MOD
        } else {
            s
        }
    }
    fn sub(x: u64, y: u64) -> u64 {
        let s = (x as i64) - (y as i64);
        if s < 0 {
            (s + MOD as i64) as u64
        } else {
            s as u64
        }
    }
    fn mul(x: u64, y: u64) -> u64 {
        ((x as u128) * (y as u128) % (MOD as u128)) as u64
    }
    fn pow(mut x: u64, mut y: u64) -> u64 {
        let mut ret = 1;
        while y != 0 {
            if y % 2 != 0 {
                ret = Self::mul(ret, x);
            }
            x = Self::mul(x, x);
            y /= 2;
        }
        ret
    }
    fn inv(x: u64) -> u64 {
        Self::pow(x, MOD - 2)
    }
    fn ntt(inp: &[u64], inv: bool) -> Vec<u64> {
        let N = inp.len();
        assert!(N.is_power_of_two());
        let deg = N.ilog2();
        let root = if inv { Self::inv(ROOT) } else { ROOT };
        eprintln!("ntt");
        let mut arr = vec![0u64; N];
        let ptr: *mut i64 = &mut arr[0] as *mut u64 as *mut i64;
        for i in 0..N {
            let j = i.reverse_bits() >> (usize::BITS - deg);
            unsafe {
                _mm_stream_si64(ptr.wrapping_add(j), inp[i] as i64);
            }
            //arr[i] = inp[j]; // not blazing fast
        }
        eprintln!("done rev");
        for layer in 0..deg {
            eprintln!("layer {}", layer);
            let sz = 1usize << layer;
            let w0 = Self::pow(root, 1 << (NTT_DEG - layer - 1));
            arr.par_chunks_mut(2 * sz)
                .map(|chunk| {
                    let mut w = 1;
                    for j in 0..sz {
                        let a = chunk[j];
                        let b = Self::mul(chunk[j + sz], w);
                        chunk[j] = Self::add(a, b);
                        chunk[j + sz] = Self::sub(a, b);
                        w = Self::mul(w, w0);
                    }
                })
                .collect::<()>();
        }
        if inv {
            let ipow = Self::pow(Self::inv(2), deg as u64);
            for x in &mut arr {
                *x = Self::mul(*x, ipow);
            }
        }
        arr
    }
}
fn polymul_ntt<const MOD: u64, const ROOT: u64>(p1: &[u64], p2: &[u64]) -> Vec<u64> {
    let mut a = Ntt::<MOD, ROOT>::ntt(p1, false);
    let b = Ntt::<MOD, ROOT>::ntt(p2, false);
    for (x, y) in zip(&mut a, b) {
        *x = ((*x as u128) * (y as u128) % (MOD as u128)) as u64;
    }
    Ntt::<MOD, ROOT>::ntt(&a, true)
}
//const TASK_MOD : u32 = 10007;
//const TASK_ROOT: u32 = 5;
const TASK_MOD: u32 = 20000159;
const TASK_ROOT: u32 = 7;
fn addmod(a: u32, b: u32) -> u32 {
    (a + b) % TASK_MOD
}
fn submod(a: u32, b: u32) -> u32 {
    (a - b) % TASK_MOD
}
fn mulmod(a: u32, b: u32) -> u32 {
    ((a as u64) * (b as u64) % (TASK_MOD as u64)) as u32
}
fn powmod(mut x: u32, mut y: u64) -> u32 {
    let mut ret = 1;
    while y != 0 {
        if y % 2 != 0 {
            ret = mulmod(ret, x);
        }
        x = mulmod(x, x);
        y /= 2;
    }
    ret
}
static mut INVERSES: [u32; TASK_MOD as usize] = [0; TASK_MOD as usize];
fn invmod(x: u32) -> u32 {
    unsafe { INVERSES[x as usize] }
}
fn polymul(p1: &[u32], p2: &[u32]) -> Vec<u32> {
    // use two moduli and CRT to convolve modulo TASK_MOD without overflow issues
    const M1: u64 = 0x7fffffffc8000001;
    const R1: u64 = 4010123396067321589;
    const M2: u64 = 0x7fffffff74000001;
    const R2: u64 = 7388713459186755282;
    let ret_size = (p1.len() + p2.len() - 1).next_power_of_two();
    let mut a: Vec<u64> = p1.iter().map(|x| *x as u64).collect();
    a.resize(ret_size, 0);
    let mut b: Vec<u64> = p2.iter().map(|x| *x as u64).collect();
    b.resize(ret_size, 0);
    let c = polymul_ntt::<M1, R1>(&a, &b);
    let d = polymul_ntt::<M2, R2>(&a, &b);
    zip(c, d)
        .map(|(x, y)| {
            // a % m1 = r1 -> a = k*m1 + r1
            // a % m2 = r2 -> (k*m1 + r1) % m2 = r2 -> k*m1 \equiv r2-r1 -> k = (r2-r1)/m1 % m2
            const INV12: u128 = 6588122875245263338; // 1/m1 is precomputed
            let k = ((((y + M2 - x) % M2) as u128 * INV12) % (M2 as u128)) as u64;
            (((k as u128) * (M1 as u128) + (x as u128)) % (TASK_MOD as u128)) as u32
        })
        .take(p1.len() + p2.len() - 1)
        .collect()
}
fn fft(arr: &[u32], root: u32) -> Vec<u32> {
    let mut rpows = vec![1u32];
    let mut cur = root;
    while cur != 1 {
        rpows.push(cur);
        cur = mulmod(cur, root);
    }
    let N = arr.len();
    let mut a = vec![0u32; N + 1];
    let mut b = vec![0u32; 2 * N];
    for i in 0..N + 1 {
        let aexp = -((N as i64) - (i as i64)) * ((N as i64) - (i as i64) - 1) / 2;
        let aexp = (aexp % (rpows.len() as i64) + (rpows.len() as i64)) % (rpows.len() as i64);
        if i != 0 {
            a[i] = mulmod(arr[N - i], rpows[aexp as usize]);
        } else {
            a[i] = 0;
        }
    }
    for i in 0..(2 * N as i64) {
        b[i as usize] = rpows[((i * (i - 1)) / 2 % (rpows.len() as i64)) as usize];
    }
    let c = polymul(&a, &b);
    let mut ret = vec![0u32; N];
    for i in 0..N {
        let p = (i as i64) * (i as i64 - 1) / 2 % (rpows.len() as i64);
        let x = mulmod(c[i + N], invmod(rpows[p as usize]));
        ret[i] = x;
    }
    ret
}
fn ifft(arr: &[u32], root: u32) -> Vec<u32> {
    let mut ret = fft(arr, invmod(root));
    let k = invmod(arr.len() as u32);
    for x in &mut ret {
        *x = mulmod(*x, k);
    }
    ret
}

fn eval_at(poly: &[u32], at: u32) -> u32 {
    let mut ret = 0u32;
    let mut t = 1u32;
    for x in poly {
        ret = addmod(ret, mulmod(t, *x));
        t = mulmod(t, at);
    }
    ret
}
fn gen() {
    const N: usize = TASK_MOD as usize;
    const K: usize = N - 3;
    let mut rng = rand::thread_rng();
    let mut poly = vec![0u32; K];
    for coef in &mut poly {
        *coef = rng.gen_range(0..TASK_MOD);
    }
    poly.resize(N, 0);
    let ft = fft(&poly, TASK_ROOT);
    let mut vals = vec![0i32; TASK_MOD as usize];
    for i in 0..TASK_MOD - 1 {
        vals[powmod(TASK_ROOT, i as u64) as usize] = ft[i as usize] as i32;
    }
    vals[0] = eval_at(&poly, 0) as i32;
    //for i in 0..N {
    //    assert!(vals[i] as u32 == eval_at(&poly, i as u32));
    //}
    vals[69] = -1;
    vals[420] = -1;
    vals[1337] = -1;
    poly.resize(K, 0);
    let flag = std::fs::read("flag.txt").unwrap();
    let mut sha = Sha512::new();
    sha.update(format!("{:?}", poly).as_bytes());
    println!(
        "{:?}",
        zip(flag, sha.finalize())
            .map(|(x, y)| x as u8 ^ y)
            .collect::<Vec<u8>>()
    );
    println!("{:?}", vals);
}
fn add(set: &[u32], d: i32) -> Vec<u32> {
    fn ad(x: u32, y: i32) -> u32 {
        (((x as i32 + y) % (TASK_MOD as i32) + (TASK_MOD as i32)) % (TASK_MOD as i32)) as u32
    }
    set.iter().map(|x| ad(*x, d)).collect()
}
fn rem(mut set: Vec<u32>, tgt: u32) -> Vec<u32> {
    if !set.contains(&tgt) {
        set.push(tgt);
    }
    set
}
static mut FACTORIALS: [u32; TASK_MOD as usize] = [0; TASK_MOD as usize];
fn factorial(i: u32) -> u32 {
    unsafe { FACTORIALS[i as usize] }
}
fn product_excl(excl: &[u32]) -> u32 {
    let mut ret = factorial(TASK_MOD - 1);
    let mut excl0 = false;
    for x in excl {
        if *x == 0 {
            excl0 = true;
            continue;
        }
        ret = mulmod(ret, invmod(*x));
    }
    if !excl0 {
        0
    } else {
        ret
    }
}
fn interp(xvals: &[u32], yvals: &[u32], at: u32) -> u32 {
    let xset = HashSet::<u32>::from_iter(xvals.iter().cloned());
    let mut excl = Vec::<u32>::new();
    for i in 0..TASK_MOD {
        if !xset.contains(&i) {
            excl.push(i);
        }
    }
    let mut ret = 0u32;
    for i in 0..xvals.len() {
        let cset = rem(excl.clone(), xvals[i]);
        let num = product_excl(&add(&cset, -(at as i32)));
        let denom = product_excl(&add(&cset, -(xvals[i] as i32)));
        ret = addmod(ret, mulmod(yvals[i], mulmod(num, invmod(denom))));
    }
    ret
}
fn parse_line(mut l: String) -> Vec<u32> {
    l.remove(0);
    l.pop();
    l.split(", ")
        .map(|x| x.parse::<i32>().unwrap())
        .map(|x| if x == -1 { TASK_MOD } else { x as u32 })
        .collect()
}
fn solve() {
    const N: usize = TASK_MOD as usize;
    let lines: Vec<String> = BufReader::new(std::io::stdin())
        .lines()
        .map(|x| x.unwrap())
        .collect();
    let encflag = parse_line(lines[0].clone());
    let mut arr = parse_line(lines[1].clone());
    let mut xs = Vec::<u32>::new();
    let mut ys = Vec::<u32>::new();
    let mut gaps = 0;
    for i in 0..N {
        if arr[i] != TASK_MOD {
            xs.push(i as u32);
            ys.push(arr[i]);
        }
    }
    for i in 0..N {
        if arr[i] == TASK_MOD {
            gaps += 1;
            arr[i] = interp(&xs, &ys, i as u32);
        }
    }
    let mut ft = vec![0u32; N - 1];
    for i in 0..(N - 1) {
        ft[i] = arr[powmod(TASK_ROOT, i as u64) as usize];
    }
    let mut ft = ifft(&ft, TASK_ROOT);
    eprintln!("ifft done\n");
    for _ in 0..(gaps - 1) {
        ft.pop().unwrap();
    }
    let mut sha = Sha512::new();
    sha.update(format!("{:?}", ft).as_bytes());
    println!(
        "{}",
        String::from_utf8(
            zip(encflag, sha.finalize())
                .map(|(x, y)| x as u8 ^ y)
                .collect()
        )
        .unwrap()
    );
}

fn main() {
    unsafe {
        FACTORIALS[0] = 1;
        for i in 1..TASK_MOD {
            FACTORIALS[i as usize] = mulmod(FACTORIALS[i as usize - 1], i);
            INVERSES[i as usize] = powmod(i, (TASK_MOD - 2) as u64);
        }
    }
    let arg = std::env::args().nth(1).unwrap();
    if arg == "gen" {
        gen();
    } else if arg == "solve" {
        solve();
    } else {
        println!("specify \"gen\" or \"solve\"");
    }
}
