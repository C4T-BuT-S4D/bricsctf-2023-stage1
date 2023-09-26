DB_NAME = "moarleeks.db"
DUMP_USER_EMAIL = "avgle4kenj0yer.kRGAkbfctDFlIA@proton.me"
DUMP_USER_USERNAME = "avgle4kenj0yer"
DUMP_USER_PASSWORD = "ysR6Z1TnBsFS5QqoXMV4ZQ"

try:
    with open("instance/secret", "rb") as f:
        SECRET_KEY = f.read()
except FileNotFoundError:
    import os

    SECRET_KEY = os.urandom(32)
    with open("instance/secret", "wb") as f:
        f.write(SECRET_KEY)
