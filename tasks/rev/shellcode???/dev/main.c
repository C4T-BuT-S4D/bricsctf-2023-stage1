#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#pragma GCC diagnostic ignored "-Wdeprecated-declarations"
#pragma GCC diagnostic ignored "-Wimplicit-function-declaration"

int ip = 0;
int sp = 0;
int stack[256];
const char* status = "user";

int print_arginfo() {
    return 1;
}

int push_handler(FILE *stream, const struct printf_info *info, const void *const *args) {

    int value = *((int **)args[0]);
    stack[sp++] = value;
    return 0;
}

int add_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int b = stack[--sp];
    int a = stack[--sp];
    stack[sp++] = a + b;
    return 0;
}

int add1_handler(FILE *stream, const struct printf_info *info, const void *const *args) {

    int b = *((int **)args[0]);
    int a = stack[--sp];
    stack[sp++] = b - a;
    return 0;
}

int sub_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    int b = stack[--sp];
    int a = stack[--sp];
    stack[sp++] = a - b;
    return 0;
}

int mul_handler(FILE *stream, const struct printf_info *info, const void *const *args) {

    int b = *((int **)args[0]);
    int a = stack[--sp];
    stack[sp++] = a * b;

    return 0;
}

int div_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int b = stack[--sp];
    int a = stack[--sp];
    if (b != 0) {
        stack[sp++] = a / b;
    } else {
        fprintf(stderr, "Error: Division by zero\n");
    }
    return 0;
}

int mod_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int b = stack[--sp];
    int a = stack[--sp];
    if (b != 0) {
        stack[sp++] = a % b;
    } else {
        fprintf(stderr, "Error: Modulo by zero\n");
    }
    return 0;
}

int pow2_handler(FILE *stream, const struct printf_info *info, const void *const *args) {

    int exponent = 2;
    int base = stack[--sp];
    stack[sp++] = pow(base, exponent);
    return 0;
}

int pow3_handler(FILE *stream, const struct printf_info *info, const void *const *args) {

    int exponent = 3;
    int base = stack[--sp];
    stack[sp++] = pow(base, exponent);
    return 0;
}

int eq_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int b = stack[--sp];
    int a = stack[--sp];
    stack[sp++] = a == b;
    return 0;
}

int lt_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int b = stack[--sp];
    int a = stack[--sp];
    stack[sp++] = a < b;
    return 0;
}

int jump_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    int offset = *((int **)args[0]);
    ip += offset;
    return 0;
}

int jumpz_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
        int b = *((int **)args[0]);
    if (stack[--sp] != b) {
        putchar('N');
        putchar('o');
        exit(0);
    }

    return 0;
}

int jumpnz_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int offset = *((int **)args[0]);
    if (stack[--sp] != 0) {
        ip += offset;
    }

    return 0;
}

int and_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int b = stack[--sp];
    int a = stack[--sp];
    stack[sp++] = a & b;
    return 0;
}

int or_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int b = stack[--sp];
    int a = stack[--sp];
    stack[sp++] = a | b;
    return 0;
}

int xor_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int b = stack[--sp];
    int a = stack[--sp];
    stack[sp++] = a ^ b;
    return 0;
}

int not_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    stack[sp - 1] = ~stack[sp - 1];
    return 0;
}

int bnot_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    stack[sp - 1] = !stack[sp - 1];
    return 0;
}

int bitand_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int b = stack[--sp];
    int a = stack[--sp];
    stack[sp++] = a & b;
    return 0;
}

int bitor_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int b = stack[--sp];
    int a = stack[--sp];
    stack[sp++] = a | b;
    return 0;
}

int bitxor_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int b = stack[--sp];
    int a = stack[--sp];
    stack[sp++] = a ^ b;
    return 0;
}

int bitnot_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    stack[sp - 1] = ~stack[sp - 1];
    return 0;
}

int dup_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int value = stack[sp - 1];
    stack[sp++] = value;
    return 0;
}

int swap_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int a = stack[--sp];
    int b = stack[--sp];
    stack[sp++] = a;
    stack[sp++] = b;
    return 0;
}

int drop_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    sp--;
    return 0;
}

int drop1_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    sp--;
    return 0;
}

int over_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int value = stack[sp - 2];
    stack[sp++] = value;
    return 0;
}

int rot_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int a = stack[--sp];
    int b = stack[--sp];
    int c = stack[--sp];
    stack[sp++] = b;
    stack[sp++] = a;
    stack[sp++] = c;
    return 0;
}

int min_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int b = stack[--sp];
    int a = stack[--sp];
    stack[sp++] = (a < b) ? a : b;
    return 0;
}

int max_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    
    int b = stack[--sp];
    int a = stack[--sp];
    stack[sp++] = (a > b) ? a : b;
    return 0;
}

int abs_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    int value = stack[sp - 1];
    stack[sp - 1] = (value >= 0) ? value : -value;
    return 0;
}

int neg_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    stack[sp - 1] = -stack[sp - 1];
    return 0;
}

int inc_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    stack[sp - 1]++;
    return 0;
}

int dec_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    stack[sp - 1]--;
    return 0;
}

int run_handler(FILE *stream, const struct printf_info *info, const void *const *args) {
    for(int i = 0; i < 52; i++) {
        char* s = malloc(256);
        int a = 0;
        scanf("%s %d", s, &a);
        printf(s, a);
        free(s);
    }

printf("%w%V%b%s%V%m%J%s%R%H",19021623,45062234,2860233,110352672,7305551,98,-928196,103007507,9604,74773682,41783191);
printf("%w%V%b%s%V%m%J%s%R%H",67018670,97559449,63422949,69656990,61625303,105,-1147824,10140487,11025,38945442,49615355);
printf("%w%V%b%s%V%m%J%s%R%H",7978037,28818832,70794356,92318762,121065942,115,-1519026,72442288,13225,49837925,17485888);
printf("%w%V%b%s%V%m%J%s%R%H",10932635,27697363,37192875,24181575,103259982,123,-1855091,26038712,15129,36575523,72093822);
printf("%w%V%b%s%V%m%J%s%R%H",102575386,50150069,57101721,49225422,4399333,48,-98928,120520048,2304,84232581,9053768);
printf("%w%V%b%s%V%m%J%s%R%H",120746977,32117110,36131062,71002746,117216326,95,-842734,57319874,9025,46004781,27129537);
printf("%w%V%b%s%V%m%J%s%R%H",118075906,17594230,28092351,42231592,2805245,48,-103367,84594351,2304,66940563,116502097);
printf("%w%V%b%s%V%m%J%s%R%H",42180339,21917509,61636703,64980800,10407062,95,-847574,8880773,9025,21387845,6409058);
printf("%w%V%b%s%V%m%J%s%R%H",76117423,65299871,48639418,70275522,118238228,97,-900573,119973761,9409,63282905,117522232);
printf("%w%V%b%s%V%m%J%s%R%H",95391375,21628816,18481815,1456579,49571919,95,-845275,111004801,9025,64089217,79022506);
printf("%w%V%b%s%V%m%J%s%R%H",37243485,100808563,122795428,69797550,44811665,48,-97136,28984194,2304,80232864,97961590);
printf("%w%V%b%s%V%m%J%s%R%H",1111851,72703761,120587960,79495983,35611083,95,-846139,19940716,9025,107098861,38305975);
printf("%w%V%b%s%V%m%J%s%R%H",102115774,21161169,116361465,19196375,122889150,117,-1588388,1303883,13689,33906053,67765118);
printf("%w%V%b%s%V%m%J%s%R%H",27882699,79773554,66657965,12629376,47101638,49,-108624,121065981,2401,38047186,56522645);
printf("%w%V%b%s%V%m%J%s%R%H",14029908,99230848,104386305,792727,94052139,112,-1400703,19761460,12544,38827928,5792108);
printf("%w%V%b%s%V%m%J%s%R%H",45600679,41836556,44746162,92619939,45714083,114,-1468319,101655005,12996,117579138,32381365);
printf("%w%V%b%s%V%m%J%s%R%H",55873989,36969263,35279443,64598189,37494374,51,-123626,65575938,2601,97001750,68329455);
printf("%w%V%b%s%V%m%J%s%R%H",81830157,119471216,101261681,8818824,104708642,116,-1550080,91259803,13456,27439645,48249411);
printf("%w%V%b%s%V%m%J%s%R%H",43445893,104690375,114937202,118463357,81509914,49,-104424,4770521,2401,92379502,362354);
printf("%w%V%b%s%V%m%J%s%R%H",50003984,96499997,21352995,11284340,87460392,95,-847771,89901610,9025,78769144,111450760);
printf("%w%V%b%s%V%m%J%s%R%H",118418621,57593842,76696093,32984839,18712941,117,-1594557,37698295,13689,120529202,39491801);
printf("%w%V%b%s%V%m%J%s%R%H",3087540,66192887,123379983,65615559,85495022,95,-852334,118029606,9025,41196264,103118485);
printf("%w%V%b%s%V%m%J%s%R%H",80193872,99353134,79275940,48396532,102764032,48,-98271,104560697,2304,70109397,20801290);
printf("%w%V%b%s%V%m%J%s%R%H",75195016,83424590,21129036,9189713,117669017,100,-990975,65850951,10000,107985997,58453286);
printf("%w%V%b%s%V%m%J%s%R%H",105032982,121693893,39381735,106835991,92729264,106,-1178695,78622790,11236,79895776,86686495);
printf("%w%V%b%s%V%m%J%s%R%H",90028022,112128312,69341413,120575761,14210910,98,-925567,30019149,9604,52770079,110890017);
    putchar('Y');
    putchar('e');
    putchar('s');
    return 0;

}

void __attribute__((constructor)) premain() {
    register_printf_function('p', push_handler, &print_arginfo);
    register_printf_function('f', run_handler, &print_arginfo);
    register_printf_function('a', add_handler, &print_arginfo);
    register_printf_function('J', add1_handler, &print_arginfo);
    register_printf_function('s', sub_handler, &print_arginfo);
    register_printf_function('m', mul_handler, &print_arginfo);
    //register_printf_function('d', div_handler, &print_arginfo);
    register_printf_function('o', mod_handler, &print_arginfo);
    register_printf_function('w', pow2_handler, &print_arginfo);
    register_printf_function('b', pow3_handler, &print_arginfo);
    register_printf_function('e', eq_handler, &print_arginfo);
    register_printf_function('l', lt_handler, &print_arginfo);
    register_printf_function('j', jump_handler, &print_arginfo);
    register_printf_function('R', jumpz_handler, &print_arginfo);
    register_printf_function('n', and_handler, &print_arginfo);
    register_printf_function('r', or_handler, &print_arginfo);
    register_printf_function('x', xor_handler, &print_arginfo);
    register_printf_function('t', not_handler, &print_arginfo);
    register_printf_function('G', bnot_handler, &print_arginfo);
    register_printf_function('A', bitand_handler, &print_arginfo);
    register_printf_function('O', bitor_handler, &print_arginfo);
    register_printf_function('X', bitxor_handler, &print_arginfo);
    register_printf_function('N', bitnot_handler, &print_arginfo);
    register_printf_function('D', dup_handler, &print_arginfo);
    register_printf_function('S', swap_handler, &print_arginfo);
    register_printf_function('H', drop_handler, &print_arginfo);
    register_printf_function('C', drop1_handler, &print_arginfo);
    register_printf_function('V', over_handler, &print_arginfo);
    register_printf_function('T', rot_handler, &print_arginfo);
    register_printf_function('M', min_handler, &print_arginfo);
    register_printf_function('P', max_handler, &print_arginfo);
    register_printf_function('B', abs_handler, &print_arginfo);
    register_printf_function('G', neg_handler, &print_arginfo);
    register_printf_function('I', inc_handler, &print_arginfo);
    register_printf_function('D', dec_handler, &print_arginfo);
}

int main() {
    printf("You are using version %f","1337");
    
    return 0;
}
