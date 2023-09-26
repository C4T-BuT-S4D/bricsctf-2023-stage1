from flask_wtf import FlaskForm
from pycountry import countries
from wtforms import (
    FileField,
    PasswordField,
    SelectField,
    StringField,
    TextAreaField,
    validators,
)


# SigninForm is the login form shown on /signin.
class SigninForm(FlaskForm):
    email = StringField(
        "Email",
        validators=[
            validators.InputRequired(message="Email required."),
            validators.Email(),
        ],
    )

    password = PasswordField(
        "Password", validators=[validators.InputRequired(message="Password required.")]
    )


# PostForm is the post creation form shown on /post.
class PostForm(FlaskForm):
    title = StringField(
        "Title",
        validators=[validators.InputRequired(message="Title required.")],
    )

    content = TextAreaField(
        "Content", validators=[validators.InputRequired(message="Content required.")]
    )

    file = FileField(
        "File", validators=[validators.InputRequired(message="File required.")]
    )


# SettingsForm is the user settings form shown on /user.
# The email and password fields aren't editable which is why they aren't validated.
class SettingsForm(FlaskForm):
    email = StringField("Email")

    password = PasswordField("Password")

    username = StringField(
        "Username", validators=[validators.InputRequired(message="Username required.")]
    )

    country = SelectField("Country", choices=[(c.alpha_2, c.name) for c in countries])
