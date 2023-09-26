import config
import forms
import repo
from flask import (
    Flask,
    flash,
    redirect,
    render_template,
    request,
    send_file,
    session,
    url_for,
)
from flask_login import (
    LoginManager,
    current_user,
    login_required,
    login_user,
    logout_user,
)

app = Flask(__name__)
app.secret_key = config.SECRET_KEY
app.config["SESSION_COOKIE_SAMESITE"] = "Strict"
app.config["SQLALCHEMY_DATABASE_URI"] = f"sqlite+pysqlite:///{config.DB_NAME}"

repo.db.init_app(app)

login_manager = LoginManager()
login_manager.init_app(app)


@login_manager.user_loader
def load_user(user_id: str):
    try:
        return repo.get_user(int(user_id))
    except ValueError:
        return None


@login_manager.unauthorized_handler
def unauthorized():
    return redirect(url_for("sign_in"))


@app.get("/")
@login_required
def index():
    return render_template(
        "index.html",
        recent=session["recent"],
        posts=repo.list_posts(),
    )


@app.route("/post", methods=["GET", "POST"])
@login_required
def create_post():
    form = forms.PostForm()

    if request.method == "POST" and form.validate():
        form.title.errors.append(
            "Post creation isn't available to everyone yet. Request access from the admin."
        )

    return render_template("create_post.html", recent=session["recent"], form=form)


@app.get("/posts/<int:post>")
@login_required
def view_post(post: int):
    post = repo.get_post(post)
    if post is None:
        flash("No such post exists", "error")
        return redirect(url_for("index"))

    recent: list[repo.Post] = session["recent"]
    recent_index = [i for i, p in enumerate(recent) if p["id"] == post.id]

    post_dict = {"id": post.id, "title": post.title}

    if len(recent_index) == 0:
        # post isn't contained in the recent list at all, push it to the top
        recent.append(post_dict)
    elif recent_index[0] == len(recent) - 1:
        # post is the last entry anyway, don't update anything
        pass
    else:
        # move post to the last position
        recent.pop(recent_index[0])
        recent.append(post_dict)

    # limit to last 3
    recent = recent[-3:]
    session["recent"] = recent

    return render_template(
        "view_post.html",
        recent=recent,
        post=post,
        unlocked=repo.check_access(current_user.id, post.id),
    )


@app.get("/files/<int:post>")
@login_required
def get_file(post: int):
    if not repo.check_access(current_user.id, post):
        flash("Access denied", "error")
        return redirect(url_for("view_post", post=post))

    post = repo.get_post(post)
    if post is None:
        flash("No such post exists", "error")
        return redirect(url_for("index"))

    return send_file(f"files/{post.id}", as_attachment=True, download_name=post.file)


@app.route("/user", methods=["GET", "POST"])
@login_required
def edit_user():
    form = forms.SettingsForm(
        username=current_user.username,
        country=current_user.country,
    )

    if request.method == "POST" and form.validate():
        repo.update_user_info(current_user, form.username.data, form.country.data)
        return redirect(url_for("edit_user"))

    return render_template(
        "edit_user.html",
        recent=session["recent"],
        user=current_user,
        form=form,
    )


@app.route("/signin", methods=["GET", "POST"])
def sign_in():
    form = forms.SigninForm()

    if request.method == "POST" and form.validate():
        user = repo.get_user_by_creds(form.email.data, form.password.data)
        if user is None:
            form.password.errors.append("Invalid credentials.")
            return render_template("signin.html", authorized=False, form=form), 401

        login_user(user)
        session["recent"] = []
        return redirect(url_for("index"))

    return render_template("signin.html", form=form)


@app.get("/signout")
@login_required
def sign_out():
    logout_user()
    return redirect(url_for("sign_in"))
