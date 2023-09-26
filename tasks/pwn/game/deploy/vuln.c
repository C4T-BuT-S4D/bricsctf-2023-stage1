#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#include <string.h>
#include <malloc.h>
#include "cards.c"


struct maze {
    int h, w;
    char *maze;
    int start_x, start_y;
    int x, y;
};
typedef struct maze maze_t;

typedef struct {
    int maze_wins, cards_wins;
    int maze_total, cards_total;
    int name_sz;
    char *name;
} user_t;

maze_t maze;
user_t *user;
char maze_backup[0x200];

char check_x_y(int x, int y, int h, int w) {
    return (0 <= x && x < h) && (0 <= y && y < w);
}

void restore_maze() {
    free(maze.maze);
    maze.maze = malloc(maze.h * maze.w+1);
    memset(maze.maze, 0, malloc_usable_size(maze.maze));
    memcpy(maze.maze, maze_backup, maze.h*maze.w);
    memset(maze_backup, 0, sizeof(maze_backup));
}

void edit_maze() {
    int h, w;

    printf("Height (3-16): ");
    scanf("%d", &h);
    if (!(3 <= h && h <= 16)) {
        printf("Invalid\n");
        return;
    }
    printf("Width (3-16): ");
    scanf("%d", &w);
    getchar();
    if (!(3 <= w && w <= 16)) {
        printf("Invalid\n");
        return;
    }

    memset(maze_backup, 0, sizeof(maze_backup));
    memcpy(maze_backup, maze.maze, maze.h*maze.w);
    free(maze.maze);

    maze.maze = malloc(h*w+1);
    memset(maze.maze, 0, malloc_usable_size(maze.maze));

    printf("Enter maze ('#' - wall, ' ' - empty): \n");
    for (int i = 0; i < h; ++i) {
        if (read(0, maze.maze+w*i, w) != w) {
            printf("error\n");
            restore_maze();
            return;
        }
        getchar();
    }
    for (int i = 0; i < w*h; ++i) {
        if (maze.maze[i] != '#' && maze.maze[i] != ' ') {
            printf("Illegal character\n");
            restore_maze();
            return;
        }
    }
    memset(maze_backup, 0, sizeof(maze_backup));

    int cnt;
    while (1) {
        printf("Portals cnt (0-5): ");
        scanf("%d", &cnt);
        if (!(0 <= cnt && cnt <= 5)) {
            printf("Invalid\n");
            continue;
        }
        break;
    }
    int i = 0;
    while (i < cnt) {
        printf("Enter portal %d (format `x1 y1 x2 y2`): ", i+1);
        int x1, x2, y1, y2;
        scanf("%d %d %d %d", &x1, &y1, &x2, &y2);
        if (!check_x_y(x1, y1, h, w) || !check_x_y(x2, y2, h, w) || (x1 == x2 && y1 == y2)) {
            printf("Invalid\n");
            continue;
        }
        if (maze.maze[x1*w+y1] != ' ' || maze.maze[x2*w+y2] != ' ') {
            printf("Invalid\n");
            continue;
        }
        maze.maze[x1*w+y1] = i+0x41;
        maze.maze[x2*w+y2] = i+0x41;
        i += 1;
    }

    maze.h = h;
    maze.w = w;

    while (1) {
        printf("Enter start pos (format `x y`): ");
        int x, y;
        scanf("%d %d", &x, &y);
        if (!check_x_y(x, y, h, w)) {
            printf("Invalid start pos\n");
            continue;
        }
        if (x == 0 || x == h-1 || y == 0 || y == w-1) {
            printf("Start point can't be on side\n");
            continue;
        }
        if (maze.maze[x*w+y] != ' ') {
            printf("Start point must be empty\n");
            continue;
        }
        maze.start_x = x;
        maze.start_y = y;
        maze.x = x;
        maze.y = y;
        maze.maze[x*w+y] = '@';
        break;
    }

}

void print_maze() {
    for (int i = 0; i < maze.h; ++i) {
        write(1, maze.maze+maze.w*i, maze.w);
        printf("\n");
    }
    printf("\nYou now at %d %d (you may not see '@' because of portal)\n", maze.x, maze.y);
}

void set_name() {
    printf("Name (%d chars max): ", user->name_sz);
    read(0, user->name, user->name_sz);
}


void user_info() {
    printf("Name: %s\nMaze wins: %d/%d\nCards wins: %d/%d\n", user->name, user->maze_wins, user->maze_total, user->cards_wins, user->cards_total);
    printf("Change name? (y/n) : ");
    char c[4];
    read(0, c, 2);
    if (c[0] == 'y') {
        set_name();
    }
}

void play_maze() {
    printf("Use a s d w keys to move, you win when you reach any maze side\n\n");
    printf("Additional keys: \ne - edit maze\ni - user info\nq - quit\n\n");

    ++user->maze_total;
    int new_x = 0, new_y = 0;

    while (1) {
        new_x = maze.x;
        new_y = maze.y;
        print_maze();
        char s[4];
        read(0, s, 2);
        char c = s[0];
        if (c == 'e') {
            edit_maze();
            continue;
        }
        if (c == 'i') {
            user_info();
            continue;
        }
        if (c == 'q') {
            return;
        }
        if (c == 'a') {
            if (maze.maze[maze.x*maze.w+maze.y-1] != '#') {
                new_y = maze.y - 1;
            }
            else continue;
        }
        else if (c == 'd') {
            if (maze.maze[maze.x*maze.w+maze.y+1] != '#') {
                new_y = maze.y + 1;
            }
            else continue;
        }
        else if (c == 'w') {
            if (maze.maze[(maze.x-1)*maze.w+maze.y] != '#') {
                new_x = maze.x - 1;
            }
            else continue;
        }
        else if (c == 's') {
            if (maze.maze[(maze.x+1)*maze.w+maze.y] != '#') {
                new_x = maze.x + 1;
            }
            else continue;
        }
        else {
            continue;
        }
        if (maze.maze[new_x*maze.w+new_y] >= 'A') {
            char f = 1;
            for (int i = 0; i < maze.h && f; ++i) {
                for (int j = 0; j < maze.w && f; ++j) {
                    if (maze.maze[i*maze.w+j] == maze.maze[new_x*maze.w+new_y] && (i != new_x || j != new_y)) {
                        new_x = i;
                        new_y = j;
                        f = 0;
                    }
                }
            }
            maze.maze[maze.x * maze.w + maze.y] = ' ';
            maze.x = new_x;
            maze.y = new_y;
            continue;
        }
        if (maze.maze[maze.x * maze.w + maze.y] < 'A') {
            maze.maze[maze.x * maze.w + maze.y] = ' ';
        }
        maze.maze[new_x * maze.w + new_y] = '@';
        maze.x = new_x;
        maze.y = new_y;
        if (maze.x == 0 || maze.x == maze.h-1 || maze.y == 0 || maze.y == maze.w-1) {
            maze.maze[maze.x*maze.w+maze.y] = ' ';
            maze.x = maze.start_x;
            maze.y = maze.start_y;
            maze.maze[maze.x*maze.w+maze.y] = '@';
            printf("You win\n");
            ++user->maze_wins;
            break;
        }
    }
}

void set_default_maze() {
    maze.h = 8;
    maze.w = 8;
    maze.maze = malloc(maze.h * maze.w);
    memcpy(maze.maze, "#########A ## A## #  # ##   #  ## #   ###      ## # #  ####### #", maze.h*maze.w);
    maze.start_x = maze.x = 3;
    maze.start_y = maze.y = 1;
    maze.maze[maze.x*maze.w+maze.y] = '@';
}

void play_cards() {
    node *player = 0;
    node *bot = 0;
    int x;
    char is_player = 0;
    node_pair *tmp;
    node *bot_choice;
    node *player_choice;
    char c[8];
    memset(c, 0, sizeof(c));

    user->cards_total += 1;

    tmp = init_state(player, bot);
    player = tmp->first;
    bot = tmp->second;
    free(tmp);

    printf("You and bot have 5 cards from 1 to 10. You and bot both choose one card. If your card is greater than bot's or you choose 1 and bot choose 10, you take both cards. Otherwise, bot takes both cards. Winner is the player who takes all cards.\n\nMore options (enter them instead of card):\n");
    printf("q - quit\nl - list your cards\ni - user info\n\nYour cards: ");
    print(player);
    printf("\n\n");
    while (1) {
        if (!player) {
            printf("Bot wins\n");
            return;
        }
        if (!bot) {
            printf("You win\n");
            user->cards_wins += 1;
            return;
        }

        tmp = cut(bot, bot->k);
        bot = tmp->first;
        bot_choice = tmp->second;
        free(tmp);

        printf("\n");
        if (!is_player) {
            printf("Bot: %d\n", bot_choice->k);
        }

        while (1) {
            printf("You: ");
            read(0, c, 2);
            if (c[0] == 'q') {
                clear(player);
                clear(bot);
                return;
            }
            if (c[0] == 'l') {
                printf("Your cards: ");
                print(player);
                printf("\n");
                printf("Bot cards: ");
                print(bot);
                printf("\n");
                continue;
            }
            if (c[0] == 'i') {
                user_info();
                continue;
            }
            x = atoi(c);
            if (!find(player, x)) {
                printf("You don't have this card\n");
                continue;
            }
            break;
        }

        tmp = cut(player, x);
        player = tmp->first;
        player_choice = tmp->second;
        free(tmp);

        if (is_player) {
            printf("Bot: %d\n", bot_choice->k);
        }

        if (((player_choice->k > bot_choice->k) && !(player_choice->k == 10 && bot_choice->k == 1)) || (player_choice->k == 1 && bot_choice->k == 10)) {
            player = insert(player, player_choice);
            player = insert(player, bot_choice);
            printf("You take cards\n");
        }
        else {
            bot = insert(bot, player_choice);
            bot = insert(bot, bot_choice);
            printf("Bot takes cards\n");
        }
    }
}


void game_menu() {
    while (1) {
        printf("1. Maze\n2. Cards\n3. User info\n4. Exit\n> ");
        int c;
        scanf("%d", &c);
        if (c == 1) {
            play_maze();
        }
        else if (c == 2) {
            play_cards();
        }
        else if (c == 3) {
            user_info();
        }
        else if (c == 4) {
            break;
        }
    }
}

int main() {
    setvbuf(stdin, NULL, _IONBF, 0);
    setvbuf(stdout, NULL, _IONBF, 0);
    user = malloc(sizeof(user_t));
    user->name_sz = 0x18;
    user->name = malloc(user->name_sz+1);
    memset(user->name, 0, malloc_usable_size(user->name));
    set_name();
    set_default_maze();
    game_menu();
    return 0;
}
