# youtube channels unsubscriber

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Unsubscribe all youtube channels subscriptions

There are 2 ways you can do this -

- [unsubscribe.js](./unsubscribe.js): A script that you run in your browser's Developer console
  - Go to <https://www.youtube.com/feed/channels>, make sure you are logged into the correct Google account
  - Review the code in `unsubscribe.js`, copy the contents and paste it in your browser's Developer console and press Enter
  - The script is a little buggy and you may need to run it a couple of times to unsubscribe from all channels, any improvements are welcome!

- `go run main.go`: This uses the Youtube API and you'll need to generate OAuth credentials to use it
  - This option is limited to 200 channels per day, see [quota usage](https://developers.google.com/youtube/v3/getting-started#quota)
