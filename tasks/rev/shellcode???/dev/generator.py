from random import randint
useless_operation = [
'%D%V%S%H',
'%D%D%V%S%H%H',
'%V%S%H%D',
'%D%V%S%H%H',
'%D%V%S%H%H',
'%H%V%S%H',
'%D%V%S%H%H',
]
def new_operation(x):
    return x
#x
#y^2 - x^3
#x*98
#%w%V%Z%s%V%m%J%s%z
# y^2 = x^3 + 98x + -928196
# PUSH X 
# PUSH Y 
# POW Y^2 
# OVER
# POW X^3
# SUB
# OVER
# mul
# ADD FROM STACK
# sub

def find_curve_equation1(x, y):
    a = x
    b = -(x**3 - y**2)
    eq = f"y^2 = x^3 + {a}x + {b}"

    ans = -(y**2 - x**3 - a*x - b)
    eq = f"y**2 - x**3 - {a}*x - {b} == {-ans}"

    return eq

def find_curve_equation(x, y):
    a = x
    b = -(x**3 - y**2)
    eq = f"y**2 == x**3 + {a}*x + {b}"
    ans = -(y**2 - x**3 - a*x - b)
    return f'printf("%w%V%b%s%V%m%J%s%R%H%H",{randint(1,127)},{randint(1,127)},{randint(1,127)},{randint(1,127)},{randint(1,127)},{a},{b},{randint(1,127)},{ans},{randint(1,127)},{randint(1,127)});'


def solution(expr):
    s = 0
    for x in range(20,127):
        for y in range(20,127):
            if eval(expr):
                s+=1
    print(s)
flag = "brics+{L0l_y0U_can_n0t_jus1_pArs3_th1s_buT_G0od_job}"[::-1]
#for i in range(0,len(flag),2):
#    print(f"%p {ord(flag[i+1])}\n%p {ord(flag[i])}")

# Generate shellcode
for i in range(0,len(flag),2):
    equation = find_curve_equation1(ord(flag[i]), ord(flag[i+1]))
    print(equation)
