import datetime

import config
from app import app
from repo import Permission, Post, User, db
from sqlalchemy import select
from sqlalchemy.sql.expression import func

# these posts will be created in the DB
db_posts = [
    Post(
        id=1,
        ts=datetime.datetime(2023, 8, 11, 13, 3),
        title="13k Combolist GMail accs",
        content="13k Google accounts collected from various phishing campaigns and malware, starting from 2017 and up until 2022. Leaving it here cause it's not that useful to me anymore.",
        file="gmail13k.zip",
        author="wedabest",
    ),
    Post(
        id=2,
        ts=datetime.datetime(2023, 8, 25, 5, 37),
        title="⭐️✨ Tele+Email+Name SCAMlist ✨⭐️ | High Quality",
        content="Fav list got it from my scamfriends. Download and enjoy",
        file="TEN-scam.tgz",
        author="h0m3invader",
    ),
    Post(
        id=3,
        ts=datetime.datetime(2023, 9, 5, 17, 49),
        title="moarleeks 1000 userpass COMBO bruteforce list",
        content="Lol, admins should probably review the security of their own website. Easiest user dump i've been able to find in a while.",
        file="moarleeks.db.bak",
        author="noobiee",
    ),
    Post(
        id=4,
        ts=datetime.datetime(2023, 9, 16, 9, 14),
        title="Leaked Russian 🇷🇺 passports",
        content="gosuslugi.ru, pochta.ru combined leaks. Be safe frens.",
        file="ruspass.txt",
        author="0megal3aks",
    ),
    Post(
        id=5,
        ts=datetime.datetime(2023, 9, 22, 11, 23),
        title="BRICS+ CTF 2023 freshly LEAKED flags 🔥🔥🔥",
        content="Fresh out the oven 🔥🔥🔥 cracked the admin infra, have fun!",
        file="flags.csv",
        author="CTFl33ker",
    ),
]


# init_db initializes all of the posts, random users, and two permissions necessary to solve the task.
# Additionally, a copy of the sqlite database is made when it only contains the random users,
# and that dump will be available to the special pcap dump leak user in one of the posts.
def init_db():
    import base64
    import os
    import random
    import shutil

    from pycountry import countries
    from random_username.generate import generate_username

    GENERATE_N_USERS = 1000

    all_countries = list(countries)
    random_country = lambda: random.choice(all_countries).alpha_2

    # Prepare random users and backup the database once it contains them
    users = []
    for _ in range(GENERATE_N_USERS):
        username = generate_username(1)[0]
        users.append(
            User(
                email=username
                + "."
                + base64.urlsafe_b64encode(os.urandom(10)).decode().strip("=")
                + "@proton.me",
                username=username,
                password=base64.urlsafe_b64encode(os.urandom(16)).decode().strip("="),
                country=random_country(),
            )
        )

    db.session.add_all(users)
    db.session.commit()

    shutil.copyfile(f"./instance/{config.DB_NAME}", f"./files/3")

    # Add posts, special user, and permissions to the DB
    leaked_user = User(
        id=1337,
        email=config.DUMP_USER_EMAIL,
        username=config.DUMP_USER_USERNAME,
        password=config.DUMP_USER_PASSWORD,
        country=random_country(),
    )

    db.session.add_all(db_posts)
    db.session.add(leaked_user)

    # permission for pcap user to access the moarleeks dump,
    # and permission for one of the generated users to access the flags
    db.session.add(Permission(user_id=leaked_user.id, post_id=db_posts[2].id))
    db.session.add(
        Permission(
            user_id=random.randint(25, GENERATE_N_USERS - 25), post_id=db_posts[4].id
        )
    )

    db.session.commit()


with app.app_context():
    db.create_all()

    num_users = db.session.execute(select(func.count(User.id))).one()[0]
    if num_users == 0:
        init_db()
