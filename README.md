# OLX-Crawler
TODO

## Installation
### 1. Install docker and docker compose
Make sure you've installed Docker and Docker compose because this project is using Docker all in order to minimise installation errors
### 2. Create a spreadsheet on Google Sheets
You can copy my spreadsheet to this: https://docs.google.com/spreadsheets/d/1a99s1jK3T-wBZI1rlpLsECdi5WFnpq8h3U6uUP8oNIs/edit?usp=sharing When you have the sheet open in your browser, the URL will look something like this: https://docs.google.com/spreadsheets/d/1-XXXXXXXXXXXXXXXXXXXSgGTwY/edit#gid=0. And in this URL, 1-XXXXXXXXXXXXXXXXXXXSgGTwY is the spreadsheet's ID and it will be different for each spreadsheet.

Save the spreadsheet ID to the. env file and change the value <b>SPREED_SHEET_ID</b>

### 3. Enable Google Sheets API in  Google developers console
You can read this tutorial in step 2 and 3 https://medium.com/swlh/how-to-read-or-modify-spreadsheets-from-google-sheets-using-node-js-6f5a672bdd37#ed85

Save the Google sheets JSON credentials API to the folder crawl-olx, Change the <b>SPREED_SHEET_AUTH</b> value with file name credentials.json

### 4. Complete the contents of the .env file
<b>MASTER_PHONE_NUMBER</b> (Telfon number you use to send WHATSAAP messages) <b>TARGET_WA_MESSAGE</b> (recipients of the message), <b>PAGE_URL_OLX</b> (the URL of the website OLX you want to monitor)

### 5. Run start.sh
Run Command ```./start.sh``` And wait for the Docker compose process to run, until it displays a QR code

If you have not had the time to scan the QR code but it has expired QR code, when run CTRL + X, and Start command ```./start.sh``

Wait 30 minutes once the bot will send a message on your WhatsApp

If successful will display the response ```[*] Waiting for messages. To exit press CTRL+C```

## Legal
This code is in no way affiliated with, authorized, maintained, sponsored or endorsed by OLX or any of its
affiliates or subsidiaries. This is an independent and unofficial software. Use at your own risk.
