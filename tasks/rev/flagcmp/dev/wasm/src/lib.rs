mod befunge {
    use std::collections::HashMap;
    use std::fmt;

    use anyhow::Ok;
    use ascii::AsciiChar;
    use lazy_static::lazy_static;
    use phf::phf_map;
    use rand::prelude::Distribution;
    use rand::rngs::SmallRng;
    use rand::{distributions, SeedableRng};

    // Command declares all the possible befunge93 commands. Additionally to the defined commands,
    // Number contains a valid befunge number from 0 to 9, while Value can contain any other ascii character
    // encountered on the playfield, e.g. letters inside a string, or a negative value fitting into i8 which
    // can be used for calculations with negative numbers.
    #[derive(PartialEq, Eq, Hash, Clone, Copy, Debug)]
    pub enum Command {
        Add,          // +
        Subtract,     // -
        Multiply,     // *
        Divide,       // /
        Modulo,       // %
        Not,          // !
        Greater,      // `
        Right,        // >
        Left,         // <
        Up,           // ^
        Down,         // v
        Random,       // ?
        HorizontalIf, // _
        VerticalIf,   // |
        StringMode,   // ""
        Dup,          // :
        Swap,         // \
        Pop,          // $
        OutputInt,    // .
        OutputChar,   // ,
        Bridge,       // #
        Get,          // g
        Put,          // p
        End,          // @
        Number(u8),
        Value(i8),
    }

    static ASCII_TO_COMMAND: phf::Map<u8, Command> = phf_map! {
        b'+' => Command::Add,
        b'-' => Command::Subtract,
        b'*' => Command::Multiply,
        b'/' => Command::Divide,
        b'%' => Command::Modulo,
        b'!' => Command::Not,
        b'`' => Command::Greater,
        b'>' => Command::Right,
        b'<' => Command::Left,
        b'^' => Command::Up,
        b'v' => Command::Down,
        b'?' => Command::Random,
        b'_' => Command::HorizontalIf,
        b'|' => Command::VerticalIf,
        b'"' => Command::StringMode,
        b':' => Command::Dup,
        b'\\' => Command::Swap,
        b'$' => Command::Pop,
        b'.' => Command::OutputInt,
        b',' => Command::OutputChar,
        b'#' => Command::Bridge,
        b'g' => Command::Get,
        b'p' => Command::Put,
        b'@' => Command::End,
    };

    lazy_static! {
        static ref COMMAND_TO_ASCII: HashMap<Command, AsciiChar> = HashMap::from_iter(
            ASCII_TO_COMMAND
                .entries()
                .map(|(k, v)| (*v, AsciiChar::from_ascii(*k).unwrap())),
        );
    }

    // TryFrom char is implemented for parsing programs.
    impl TryFrom<char> for Command {
        type Error = anyhow::Error;

        fn try_from(c: char) -> Result<Self, Self::Error> {
            let ch = AsciiChar::from_ascii(c)?;

            let cmd = if ch.is_digit(10) {
                Self::Number(ch.as_byte() - AsciiChar::_0.as_byte())
            } else if let Some(cmd) = ASCII_TO_COMMAND.get(&ch.into()) {
                *cmd
            } else {
                Self::Value(ch.as_byte().try_into().unwrap())
            };

            Ok(cmd)
        }
    }

    // TryFrom i32 is implemented for converting stack values into program commands.
    impl TryFrom<i32> for Command {
        type Error = anyhow::Error;

        fn try_from(value: i32) -> Result<Self, Self::Error> {
            if value < 0 {
                Ok(Self::Value(i8::try_from(value)?))
            } else {
                Self::try_from(char::try_from(u8::try_from(value)?)?)
            }
        }
    }

    // From command is implemented for converting program commands into stack values.
    impl From<Command> for i32 {
        fn from(value: Command) -> Self {
            match value {
                Command::Number(n) => (n + AsciiChar::_0.as_byte()).into(),
                Command::Value(v) => v.into(),
                _ => COMMAND_TO_ASCII.get(&value).unwrap().as_byte().into(),
            }
        }
    }

    // Befunge-93 page size constants from the original C interpreter
    pub const LINE_WIDTH: usize = 80;
    pub const PAGE_HEIGHT: usize = 25;

    // Playfield is indexed y-first, x-second
    type Playfield = [[Command; LINE_WIDTH]; PAGE_HEIGHT];

    pub fn new_playfield() -> Playfield {
        return [[Command::Value(0); LINE_WIDTH]; PAGE_HEIGHT];
    }

    type Coords = (i8, i8);

    const DELTA_RIGHT: Coords = (1, 0);
    const DELTA_LEFT: Coords = (-1, 0);
    const DELTA_UP: Coords = (0, -1);
    const DELTA_DOWN: Coords = (0, 1);

    static DELTAS: [Coords; 4] = [DELTA_RIGHT, DELTA_LEFT, DELTA_UP, DELTA_DOWN];

    // InterpretationMode specifies how the interpreter handles the instructions on the playfield.
    enum InterpretationMode {
        Exec,
        String,
        Bridge,
    }

    // Interpreter is a step-by-step befunge-93 interpreter.
    pub struct Interpreter<W: fmt::Write> {
        mode: InterpretationMode,
        playfield: Box<Playfield>,
        stack: Vec<i32>,
        pc: Coords,    // x, y
        delta: Coords, // x, y
        rng: SmallRng,
        output: W,
    }

    impl<W: fmt::Write> Interpreter<W> {
        pub fn new(playfield: &Playfield, output: W) -> Self {
            Self {
                mode: InterpretationMode::Exec,
                playfield: Box::new(*playfield),
                stack: Vec::new(),
                pc: (0, 0),
                delta: DELTA_RIGHT,
                rng: SmallRng::from_entropy(),
                output,
            }
        }

        pub fn step(self: &mut Interpreter<W>) -> Option<()> {
            let ins = self.playfield[self.pc.1 as usize][self.pc.0 as usize];
            match self.mode {
                InterpretationMode::Exec => self.exec(ins)?,
                InterpretationMode::String => {
                    if ins == Command::StringMode {
                        self.mode = InterpretationMode::Exec
                    } else {
                        self.push(ins.into())
                    }
                }
                InterpretationMode::Bridge => self.mode = InterpretationMode::Exec,
            };

            self.pc = (
                (self.pc.0 + self.delta.0).rem_euclid(LINE_WIDTH as i8),
                (self.pc.1 + self.delta.1).rem_euclid(PAGE_HEIGHT as i8),
            );

            Some(())
        }

        fn exec(self: &mut Interpreter<W>, cmd: Command) -> Option<()> {
            match cmd {
                Command::Add => {
                    let (b, a) = (self.pop(), self.pop());
                    self.push(a.wrapping_add(b));
                }
                Command::Subtract => {
                    let (b, a) = (self.pop(), self.pop());
                    self.push(a.wrapping_sub(b));
                }
                Command::Multiply => {
                    let (b, a) = (self.pop(), self.pop());
                    self.push(a.wrapping_mul(b));
                }
                Command::Divide => {
                    let (b, a) = (self.pop(), self.pop());
                    if b == 0 {
                        return None;
                    }
                    self.push(a.wrapping_div(b));
                }
                Command::Modulo => {
                    let (b, a) = (self.pop(), self.pop());
                    if b == 0 {
                        return None;
                    }
                    self.push(a.wrapping_rem(b));
                }
                Command::Not => {
                    let a = self.pop();
                    self.push(!(a != 0) as i32);
                }
                Command::Greater => {
                    let (b, a) = (self.pop(), self.pop());
                    self.push((a > b) as i32);
                }
                Command::Right => self.delta = DELTA_RIGHT,
                Command::Left => self.delta = DELTA_LEFT,
                Command::Up => self.delta = DELTA_UP,
                Command::Down => self.delta = DELTA_DOWN,
                Command::Random => {
                    let dist = distributions::Uniform::new(0, 4);
                    self.delta = DELTAS[dist.sample(&mut self.rng)]
                }
                Command::HorizontalIf => {
                    let a = self.pop();
                    self.delta = if a != 0 { DELTA_LEFT } else { DELTA_RIGHT }
                }
                Command::VerticalIf => {
                    let a = self.pop();
                    self.delta = if a != 0 { DELTA_UP } else { DELTA_DOWN }
                }
                Command::StringMode => self.mode = InterpretationMode::String,
                Command::Dup => {
                    let a = self.pop();
                    self.push(a);
                    self.push(a);
                }
                Command::Swap => {
                    let (b, a) = (self.pop(), self.pop());
                    self.push(a);
                    self.push(b);
                }
                Command::Pop => {
                    _ = self.pop();
                }
                Command::OutputInt => {
                    let a = self.pop();
                    write!(self.output, "{}", a).ok()?;
                }
                Command::OutputChar => {
                    let a = self.pop();
                    let ch = AsciiChar::from_ascii(u8::try_from(a).ok()?).ok()?;
                    self.output.write_char(ch.as_char()).ok()?;
                }
                Command::Bridge => {
                    self.mode = InterpretationMode::Bridge;
                }
                Command::Get => {
                    let (b, a) = (self.pop(), self.pop());
                    let (x, y) = (usize::try_from(a).ok()?, usize::try_from(b).ok()?);
                    let ins = *self.playfield.get(y)?.get(x)?;
                    self.push(ins.into());
                }
                Command::Put => {
                    let (b, a) = (self.pop(), self.pop());
                    let value = self.pop();
                    let (x, y) = (usize::try_from(a).ok()?, usize::try_from(b).ok()?);
                    let pos = self.playfield.get_mut(y)?.get_mut(x)?;
                    *pos = value.try_into().ok()?;
                }
                Command::End => return None,
                Command::Number(n) => self.push(n.into()),
                Command::Value(v) => {
                    if v != AsciiChar::Space.as_byte() as i8 {
                        return None;
                    }
                }
            };

            Some(())
        }

        fn pop(self: &mut Interpreter<W>) -> i32 {
            self.stack.pop().unwrap_or(0)
        }

        fn push(self: &mut Interpreter<W>, value: i32) {
            self.stack.push(value)
        }
    }
}

use std::collections::HashSet;
use std::iter;
use std::ops;

use befunge::Command;

use lazy_static::lazy_static;
use wasm_bindgen::prelude::*;

const ITERATION_LIMIT: usize = 200000;

fn count_repr_command(repr: &Vec<Command>, cmd: Command) -> usize {
    repr.iter().filter(|&&v| v == cmd).count()
}

fn count_repr_unique(repr: &Vec<Command>) -> usize {
    let mut set = HashSet::new();
    for cmd in repr {
        set.insert(cmd);
    }
    set.len()
}

lazy_static! {
    static ref NUMBER_REPR: [Vec<Command>; 256] = {
        // allocate with more capacity so that all possible befunge operators are used (including div)
        let mut number_repr = vec![vec![]; 4000];

        let repr_ops: [(Command, fn(i32, i32) -> i32); 4] = [
            (Command::Add, ops::Add::add),
            (Command::Subtract, ops::Sub::sub),
            (Command::Multiply, ops::Mul::mul),
            (Command::Divide, ops::Div::div),
        ];

        for n in 0..number_repr.len() {
            if n < 10 {
                number_repr[n] = vec![Command::Number(n as u8)];
            }

            for other in (0..10).chain(iter::once(n)) {
                for op in repr_ops {
                    if other == 0 && op.0 == Command::Divide {
                        continue;
                    }

                    let res = op.1(n as i32, other as i32);
                    if res < 0 || res >= number_repr.len() as i32 {
                        continue;
                    }

                    let res = res as usize;

                    let other_repr = if other == n {
                        Command::Dup
                    } else {
                        Command::Number(other as u8)
                    };

                    let mut res_repr;
                    if (op.0 == Command::Add || op.0 == Command::Multiply)
                        && *number_repr[n].last().unwrap() == op.0
                        && other_repr != Command::Dup
                    {
                        // simplify x+y+ to xy++ and x*y* to xy**
                        res_repr = number_repr[n].split_last().unwrap().1.to_vec();
                        res_repr.extend([other_repr, op.0, op.0].iter());
                    } else {
                        res_repr = number_repr[n].clone();
                        res_repr.extend([other_repr, op.0].iter());
                    }

                    // set result representation if one doesn't yet exist, or if this one is more "interesting"
                    let cur_repr = &number_repr[res];
                    if cur_repr.is_empty()
                        || res_repr.len() < cur_repr.len()
                        || (res_repr.len() == cur_repr.len()
                            && (count_repr_command(&res_repr, Command::Dup)
                                > count_repr_command(cur_repr, Command::Dup)
                                || count_repr_command(&res_repr, Command::Divide)
                                    > count_repr_command(cur_repr, Command::Divide)
                                || count_repr_unique(&res_repr) > count_repr_unique(cur_repr)))
                    {
                        number_repr[res] = res_repr;
                    }
                }
            }
        }

        number_repr[..256].to_vec().try_into().unwrap()
    };
}

pub static mut CHECKER_CODE: Vec<Command> = Vec::new();

#[wasm_bindgen]
pub fn set_flag(flag: &str) {
    // set original validity variable to 1
    let mut code = vec![
        Command::Number(1),
        Command::Number(0),
        Command::Dup,
        Command::Put,
    ];

    for b in flag.bytes().rev() {
        // calculate number representing next byte
        code.extend(NUMBER_REPR[b as usize].iter());
        // check this byte and write the result
        code.extend(
            vec![
                Command::Subtract,
                Command::Not,
                Command::Number(0),
                Command::Dup,
                Command::Get,
                Command::Multiply,
                Command::Number(0),
                Command::Dup,
                Command::Put,
            ]
            .iter(),
        );
    }

    unsafe {
        CHECKER_CODE = code;
    }
}

#[wasm_bindgen]
pub fn run(guess: &str) -> Option<String> {
    // Parse to allow befunge injection (lol)
    let guess_chars = guess
        .chars()
        .map(Command::try_from)
        .collect::<Result<Vec<_>, _>>()
        .ok()?;

    // Prepare code which checks the guess and prints the result
    let mut code = Vec::new();

    // write bytes of the guess to the stack
    code.push(Command::StringMode);
    code.extend(guess_chars);
    code.push(Command::StringMode);

    // check the guess
    unsafe {
        code.extend(&CHECKER_CODE);
    }

    // print the result
    code.push(Command::Number(1));
    code.push(Command::OutputInt);
    code.extend(&NUMBER_REPR[10]);
    code.push(Command::OutputChar);
    code.extend(vec![
        Command::Number(0),
        Command::Dup,
        Command::Get,
        Command::OutputInt,
    ]);
    code.push(Command::End);

    // Format the code on the playfield
    let mut playfield = befunge::new_playfield();

    let mut line = 0;
    let mut line_len = 0;
    let mut line_i = 0;
    let mut line_iter = 1i64;
    let mut reverse = false;
    for c in code {
        if line_len == befunge::LINE_WIDTH - 1 {
            playfield[line][line_i] = Command::Down;

            line += 1;
            if line == befunge::PAGE_HEIGHT {
                // too long, unlucky
                return None;
            }

            reverse = !reverse;
            line_len = 0;

            if reverse {
                line_i = befunge::LINE_WIDTH - 1;
                line_iter = -1;
                playfield[line][line_i] = Command::Left;
            } else {
                line_i = 0;
                line_iter = 1;
                playfield[line][line_i] = Command::Right;
            };

            line_len += 1;
            line_i = ((line_i as i64) + line_iter) as usize;
        }

        playfield[line][line_i] = c;
        line_len += 1;
        line_i = ((line_i as i64) + line_iter) as usize;
    }

    let mut output = String::new();
    let mut interpreter = befunge::Interpreter::new(&playfield, &mut output);

    for _ in 0..ITERATION_LIMIT {
        if interpreter.step() == None {
            break;
        }
    }

    // read output
    let length_delimiter = output.find('\n')?;
    let output_length = output[..length_delimiter].parse::<usize>().ok()?;

    let result = &output[length_delimiter + 1..];
    if result.len() != output_length {
        return None;
    }

    Some(result.to_owned())
}
