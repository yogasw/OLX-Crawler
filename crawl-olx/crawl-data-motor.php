<?php
require("JSON.php");

header('Content-type: application/json; charset=UTF-8');
$json = new Services_JSON();
$cookieFile = "cookie.txt";
$URL_AUTH = "https://www.olx.co.id/api/auth/authenticate";
$URL_MOTOR = "https://www.olx.co.id/yogyakarta-kota_g4000072/q-macbook-pro-?sorting=desc-creation";
$URL_SCHEMA = "https://www.olx.co.id/item/";
$EMAIL = "mryoga1995@gmail.com";
$PASSWORD = "Yoga12345";

if (!file_exists($cookieFile)) {
    $fh = fopen($cookieFile, "w");
    fwrite($fh, "");
    fclose($fh);
}
getCURL("https://www.olx.co.id/");
//login();

crawlMotor();

function getCURL($url, $data = null, $header = null)
{
    print_r("Crawl url : " . $url . "\n");
    global $cookieFile;
    $ch = curl_init();
    $headers = [
        'user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.92 Safari/537.36',
        'content-type: application/json',
        'Referer: https://www.olx.co.id/',
    ];

    if ($header) array_push($headers, $header);
    $method = "GET";
    if ($data) $method = "POST";
    curl_setopt_array($ch, array(
        CURLOPT_URL => $url,
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_ENCODING => "",
        CURLOPT_MAXREDIRS => 10,
        CURLOPT_TIMEOUT => 0,
        CURLOPT_COOKIEFILE => dirname(__FILE__) . '/' . $cookieFile,
        CURLOPT_COOKIEJAR => dirname(__FILE__) . '/' . $cookieFile,
        CURLOPT_FOLLOWLOCATION => true,
        CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
        CURLOPT_CUSTOMREQUEST => $method,
        CURLOPT_POSTFIELDS => $data,
        CURLOPT_HTTPHEADER => $headers,
    ));
    $response = curl_exec($ch);
    curl_close($ch);
    return $response;
}

function login()
{
    global $EMAIL, $PASSWORD, $URL_AUTH;
    $auth_data = [
        "grantType" => "email",
        "email" => $EMAIL,
        "language" => "id"
    ];
    $auth_result = getCURL($URL_AUTH, json_encode($auth_data));
    $auth_json = json_decode($auth_result, true);

    $login_authorization = 'authorization: Bearer ' . $auth_json["token"];
    $login_data = [
        "grantType" => "password",
        "password" => $PASSWORD,
        "language" => "id"
    ];
    getCURL($URL_AUTH, json_encode($login_data), $login_authorization);
}

function crawlMotor()
{
    global $URL_MOTOR, $json;
    $data = getDataJSOLX(getCURL($URL_MOTOR));
    $data = $json->decode($data);
    $data = $data->states->items->collections;
    $data = get_object_vars($data);
    $data = reset($data);
    foreach ($data as $item) {
        crawlDetail($item);
    }
}

function crawlDetail($item)
{
    global $json, $URL_SCHEMA;
    $url = $URL_SCHEMA . $item;
    $html = getDataJSOLX(getCURL($url));
    $html = $json->decode($html);

    $elements = $html->states->items->elements->{$item};
    $id = $elements->id;
    $url = $html->props->location->hostname . $html->props->location->pathname;
    $created_at = $elements->created_at;
    $title = $elements->title;
    $description = $elements->description;
    $location = $elements->locations[0];
    $location = "https://maps.google.com/?q=" . $location->lat . "," . $location->lon;
    $price = $elements->price->value->display;
    $image = $elements->images[0]->url;
    $parameters = [];
    foreach ($elements->parameters as $item) {
        array_push($parameters, $item->key_name . " : " . $item->formatted_value);
    };
    $parameters = implode(", ", $parameters);
    print_r($title);
    print_r($parameters);
}

function getDataJSOLX($html)
{
    $html = trim(preg_replace('/[\t\n\r\s]+/', ' ', $html));
    $html = explode('<script type="text/javascript">', $html);
    if (!$html[3]) print_r($html);
    $html = explode("; </script>", $html[3]);
    $html = explode('window.__APP =', $html[0]);
    return $html[1];
}

function urlClean($url)
{
    return urldecode(stripslashes($url));
}

unlink('cookie.txt');
