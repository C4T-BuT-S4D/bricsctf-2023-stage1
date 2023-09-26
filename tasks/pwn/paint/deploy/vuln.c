#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

typedef struct {
    unsigned char width;
    unsigned char height;
    char *canvas;
    int rate;
    char *comment;
} Canvas;

Canvas *canvas[4];

int read_idx() {
    int idx;
    printf("Enter idx: ");
    scanf("%d", &idx);
    return idx;
}

void add_canvas() {
    int idx = read_idx();
    if (canvas[idx]) {
        puts("Invalid index");
        return;
    }
    canvas[idx] = malloc(sizeof(Canvas));
    unsigned int width = 0, height = 0;
    printf("Enter canvas width (1-255): ");
    scanf("%hhd", &width);
    printf("Enter canvas height (1-255): ");
    scanf("%hhd", &height);
    if (width * height >= 0x100) {
        puts("Too big");
        free(canvas[idx]);
        canvas[idx] = 0;
        return;
    }
    canvas[idx]->canvas = malloc(width * height + 1);
    canvas[idx]->width = width;
    canvas[idx]->height = height;
    canvas[idx]->rate = 0;
    canvas[idx]->comment = 0;
    puts("Done");
}

void resize_canvas() {
    int idx = read_idx();
    if (!canvas[idx]) {
        puts("Invalid index");
        return;
    }
    unsigned int width = 0, height = 0;
    printf("Enter new width (1-255): ");
    scanf("%hhd", &width);
    printf("Enter new height (1-255): ");
    scanf("%hhd", &height);
    char *new_canvas = malloc(width * height + 1);
    if (!new_canvas) {
        puts("malloc error");
        return;
    }
    canvas[idx]->canvas = new_canvas;
    canvas[idx]->width = width;
    canvas[idx]->height = height;
    puts("Done");
}

void draw() {
    int idx = read_idx();
    if (!canvas[idx]) {
        puts("Invalid index");
        return;
    }
    puts("Enter your picture (`width` chars in `height` lines): ");
    for (int i = 0; i < canvas[idx]->height; ++i) {
        read(0, &canvas[idx]->canvas[canvas[idx]->width * i], canvas[idx]->width+1);
    }
}

void show() {
    int idx = read_idx();
    if (!canvas[idx]) {
        puts("Invalid index");
        return;
    }
    puts("Picture: ");
    for (int i = 0; i < canvas[idx]->height; ++i) {
        write(1, &canvas[idx]->canvas[canvas[idx]->width * i], canvas[idx]->width);
        puts("");
    }
    printf("Rate: %d\n", canvas[idx]->rate);
    if (canvas[idx]->comment) {
        printf("Comment: %s\n", canvas[idx]->comment);
    }
}

void rate() {
    int rate, idx = read_idx();
    if (!canvas[idx]) {
        puts("Invalid index");
        return;
    }
    printf("Enter rate: ");
    scanf("%d", &rate);
    canvas[idx]->rate = rate;
    printf("Leave comment (y/n): ");
    char c;
    getchar();
    scanf("%c", &c);
    getchar();
    if (c != 'y' && c != 'n') {
        puts("Invalid choice");
        return;
    }
    if (c == 'n') return;
    char comment[0x40];
    memset(comment, 0, sizeof(comment));
    printf("Enter your comment: ");
    read(0, comment, sizeof(comment)-1);
    if (canvas[idx]->comment) {
        free(canvas[idx]->comment);
    }
    canvas[idx]->comment = strdup(comment);
    puts("Done");
}

void delete_canvas() {
    int idx = read_idx();
    if (!canvas[idx]) {
        puts("Invalid index");
        return;
    }
    if (canvas[idx]->comment) {
        free(canvas[idx]->comment);
    }
    free(canvas[idx]);
    canvas[idx] = 0;
    puts("Done");
}

void menu() {
    printf("\n1. Add canvas\n2. Resize canvas\n3. Draw\n4. Show\n5. Rate\n6. Delete canvas\n\n> ");
}

int main() {
    setvbuf(stdin, NULL, _IONBF, 0);
    setvbuf(stdout, NULL, _IONBF, 0);
    int c;
    while (1) {
        menu();
        scanf("%d", &c);
        switch (c) {
        case 1:
            add_canvas();
            break;
        case 2:
            resize_canvas();
            break;
        case 3:
            draw();
            break;
        case 4:
            show();
            break;
        case 5:
            rate();
            break;
        case 6:
            delete_canvas();
            break;
        default:
            puts("Invalid");
            return 0;
        }
    }

    return 0;
}
