# OLX-Crawler
OLX Crawler is a bot that is used to monitor OLX sites periodically and send automatically to WhatsApp when there is new item on OLX site in detail
<br>Support :  
* OLX Indonesia (https://www.olx.co.id)
* OLX India (https://www.olx.in)
* OLX Pakistan (https://www.olx.com.pk)
* OLX South Africa (https://www.olx.co.za)

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
## Example receive messages

<img src="https://raw.githubusercontent.com/arioki1/OLX-Crawler/master/extras/WhatsApp%20Image%202020-06-01%20at%2013.00.46.jpeg" width="355" title="hover text">

<b>Beat Street 2019 akhir AB sleman (km 9rb'n.pjk panjang)</b>

<b>Ad description</b><br>
<b>Created At :</b><br>
2020-06-01T12:48:28+07:00

<b>Location :</b><br>
<i>https://maps.google.com/?q=-7.803,110.438</i>

<b>Image :</b><br>
<i>https://apollo-singapore.akamaized.net:443/v1/files/s9jg9ezk972s1-ID/image,</i>
<i>https://apollo-singapore.akamaized.net:443/v1/files/x2uhjpq6ysyo2-ID/image,</i>
<i>https://apollo-singapore.akamaized.net:443/v1/files/onfoxqc25k5i1-ID/image</i>

<b>Price :</b> Rp 12.500.000

<b>Parameters : </b><br>
Merek : Honda<br>
Model : Beat<br>
Tahun : 2019<br>
Jarak tempuh : 20.000-25.000 km<br>
Tipe Penjual : Individu<br>
phone : +62813927xxxxx<br>
wa : <i>https://wa.me/62813xxx</i><br>

<b>Description : </b><br>
Jual cepat BU bgt beat street 2019 AB sleman
Standar.mulus.orisinil (km 9rb'n)
Pajak hidup bln november 2020..Stnk.bpkb & faktur + kunci serep (lengkap)
Harga:12,5_passs/nettt (nego NO RESPON)
Lokasi:Berbah (sleman)
Yg serius lgsg tlp/wa aja..
Nb:Shok blkg pake punya yamaha X_ride ori empuk & lbh tinggi (yg ori lupa nyimpen)

<b>Url : </b><br>
<i>www.olx.co.id/item/beat-street-2019-akhir-ab-sleman-km-9rbnpjk-panjang-iid-784397648</i>

<b>User description : </b><br>
<b>Name :</b> Eddy Susanto<br>

<b>Profile :</b><br>
<i>https://www.olx.co.id/profile/65695098<i/>

<b>Location :</b><br>
<i>https://maps.google.com/?q=-7.803,110.438</i>

<b>Badges :</b> Facebook, G+, Phone number

## Legal
This code is in no way affiliated with, authorized, maintained, sponsored or endorsed by OLX or any of its
affiliates or subsidiaries. This is an independent and unofficial software. Use at your own risk.
