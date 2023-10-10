# crossstitch-auditor
User self-promotion auditor script for reddit.com/r/CrossStitch

## Background
This tool is intended for use by moderators of a subreddit. When designers and shop owners share their art on Reddit, they often post to any relevant subreddits they can find. This makes it difficult for moderators of one specific subreddit to audit the user's activity in that sub to determine if the user is violating self-promotion policies (ex: Reddit's 90/10 policy).

Note: The tool is configured to run against the r/CrossStitch subreddit by default, but the exact subreddit is configurable.

## Set up

1. You need a Reddit user account to authenticate with the Reddit API. Preferably, the account you use is clearly marked as a bot account. For example, our bot account is [u/CrossStitchBot](https://www.reddit.com/user/CrossStitchBot).
1. Register the auditor application with Reddit in order to get a Reddit client secret (necessary for authenticating with Reddit). Have the client secret handy.
1. Have the Reddit password for your bot account handy.


## Using the auditor locally

1. If you are using a Linux or MacOS machine, open your terminal.
1. If you are using a Windows machine, please install [GitBash](https://gitforwindows.org/) to give you a usable terminal.
1. Download the latest release of the auditor application from GitHub, [here](https://github.com/khipkin/crossstitch-auditor/releases).
2. In your terminal, navigate to the directory where you downloaded the auditor app.
1. Substituting the values for your real Reddit client secret and Reddit password, as well as the username of the Reddit user you'd like to audit (WITHOUT the "u/"), run this command in your terminal:
   ```
    REDDIT_CLIENT_SECRET=<client_secret_here> REDDIT_PASSWORD=<password_here> REDDIT_USER=<user_to_audit_here> ./crossstitch-auditor.exe
   ```
1. If everything has been set up correctly, you will see the results of your audit as output in MD format. Copy and paste the output into your favorite Markdown editor, and do your thing!

## Hosting the auditor in the cloud

TODO
