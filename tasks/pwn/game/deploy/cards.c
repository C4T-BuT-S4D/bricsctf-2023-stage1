#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <malloc.h>


struct node {
    struct node *l, *r;
    int k, p;
};
typedef struct node node;

node *new_node(int k) {
    node *v = (node *)malloc(sizeof(node));
    v->l = 0;
    v->r = 0;
    v->k = k;
    v->p = rand();
    return v;
}

struct node_pair {
    node *first;
    node *second;
};
typedef struct node_pair node_pair;

node_pair *split(node *v, int k) {
    node_pair *res = (node_pair *)malloc(sizeof(node_pair));
    res->first = 0;
    res->second = 0;
    if (!v) return res;
    if (k <= v->k) {
        node_pair *x = split(v->l, k);
        v->l = x->second;
        res->first = x->first;
        res->second = v;
        free(x);
        return res;
    }
    else {
        node_pair *x = split(v->r, k);
        v->r = x->first;
        res->first = v;
        res->second = x->second;
        free(x);
        return res;
    }
}

node *merge(node *l, node *r) {
    if (!l) return r;
    if (!r) return l;
    if (l->p < r->p) {
        r->l = merge(l, r->l);
        return r;
    }
    else {
        l->r = merge(l->r, r);
        return l;
    }
}

void print(node *v) {
    if (!v) return;
    print(v->l);
    printf("%d ", v->k);
    print(v->r);
}

node *insert(node *root, node *v) {
    node_pair *x = split(root, v->k+1);
    node *res = merge(x->first, merge(v, x->second));
    free(x);
    return res;
}

char find(node *v, int k) {
    if (!v) return 0;
    if (v->k == k) return 1;
    if (k < v->k) {
        return find(v->l, k);
    }
    else {
        return find(v->r, k);
    }
}

node_pair *cut(node *v, int k) {
    node_pair *p1 = split(v, k);
    node_pair *p2 = split(p1->second, k+1);
    node_pair *res = (node_pair *)malloc(sizeof(node_pair));
    res->first = merge(p1->first, p2->second);
    res->second = p2->first;
    free(p1);
    free(p2);
    return res;
}

void clear(node *v) {
    if (!v) return;
    clear(v->l);
    clear(v->r);
    free(v);
}

node_pair *init_state(node *user, node *bot) {
    srand(time(0));
    node *v;
    int cnt = 0;
    int N = 10;
    int u[10] = {};
    while (cnt < (N/2)) {
        int x = rand() % (N);
        if (u[x]) continue;
        cnt += 1;
        u[x] = 1;

        v = new_node(x+1); 
        user = insert(user, v);
    }
    for (int i = 0; i < N; ++i) {
        if (!u[i]) {
            v = new_node(i+1); 
            bot = insert(bot, v);
        }
    }
    node_pair *res = (node_pair *)malloc(sizeof(node_pair));
    res->first = user;
    res->second = bot;
    return res;
}

