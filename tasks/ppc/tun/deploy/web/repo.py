import datetime
import operator
from typing import Optional

from flask_login import UserMixin
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy import DateTime, ForeignKey, select
from sqlalchemy.orm import Mapped, mapped_column


db = SQLAlchemy()


# User model intentionally contains plaintext password for the DB leak.
class User(db.Model, UserMixin):
    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)
    email: Mapped[str] = mapped_column(unique=True)
    username: Mapped[str]
    password: Mapped[str]
    country: Mapped[str]


# get_user gets a user with the specified ID, or None, if no such user exists.
def get_user(id: int) -> Optional[User]:
    result = db.session.execute(select(User).where(User.id == id)).one_or_none()
    if result is None:
        return None
    return result[0]


# get_user_by_creds returns a user with the specified credentials, or None, if no such user exists.
def get_user_by_creds(email: str, password: str) -> Optional[User]:
    result = db.session.execute(
        select(User).where(User.email == email).where(User.password == password)
    ).one_or_none()
    if result is None:
        return None
    return result[0]


# update_user_info updates the editable information of a user.
def update_user_info(user: User, username: str, country: str):
    user.username = username
    user.country = country
    db.session.commit()


# Post model doesn't contain relationship to user via author because specific preset names are used.
class Post(db.Model):
    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)
    ts: Mapped[datetime.datetime] = mapped_column(DateTime())
    title: Mapped[str]
    content: Mapped[str]
    file: Mapped[str]
    author: Mapped[str]


# list_posts lists all the existing posts
def list_posts() -> list[Post]:
    return list(map(operator.itemgetter(0), db.session.execute(select(Post)).all()))


# get_post gets a post with the specified ID, or None, if no such post exists.
def get_post(id: int) -> Optional[Post]:
    result = db.session.execute(select(Post).where(Post.id == id)).one_or_none()
    if result is None:
        return None
    return result[0]


class Permission(db.Model):
    user_id: Mapped[int] = mapped_column(ForeignKey("user.id"), primary_key=True)
    post_id: Mapped[int] = mapped_column(ForeignKey("post.id"), primary_key=True)


# check_access checks if the user with the given id has access to the specified post, returning True if so, and False if not.
def check_access(user_id: int, post_id: int) -> bool:
    result = db.session.execute(
        select(Permission)
        .where(Permission.user_id == user_id)
        .where(Permission.post_id == post_id)
    ).one_or_none()
    return result is not None
