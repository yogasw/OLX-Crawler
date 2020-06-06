<?php
require __DIR__ . "/JSON.php";
require __DIR__ . "/SpreedSheet.php";

class CrawlDataOLX
{
    public $json,
        $cookieFile,
        $EMAIL,
        $PASSWORD,
        $URL_TARGET_CRAWL,
        $parse,
        $BASE_URL,
        $URL_AUTH,
        $spreedSheet;

    public function __construct()
    {
        $this->json = new Services_JSON();
        $this->cookieFile = "cookie.txt";
        $this->EMAIL = "";
        $this->PASSWORD = "";
        $this->URL_TARGET_CRAWL = getenv("PAGE_URL_OLX");
        $this->parse = parse_url($this->URL_TARGET_CRAWL);

        $parse = $this->parse;
        $this->BASE_URL = "$parse[scheme]://$parse[host]";
        $this->URL_AUTH = "$this->BASE_URL/api/auth/authenticate";
        $this->spreedSheet = new SpreedSheet();
        if (!file_exists($this->cookieFile)) {
            $fh = fopen($this->cookieFile, "w");
            fwrite($fh, "");
            fclose($fh);
        }

        $this->crawlListData();
        //login();
        unlink('cookie.txt');
    }

    function getCURL($url, $data = null, $header = null)
    {
        print_r("Crawl url : " . $url . "\n");
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
            CURLOPT_COOKIEFILE => dirname(__FILE__) . '/' . $this->cookieFile,
            CURLOPT_COOKIEJAR => dirname(__FILE__) . '/' . $this->cookieFile,
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
        $auth_data = [
            "grantType" => "email",
            "email" => $this->EMAIL,
            "language" => "id"
        ];
        $auth_result = $this->getCURL($this->URL_AUTH, json_encode($auth_data));
        $auth_json = json_decode($auth_result, true);

        $login_authorization = 'authorization: Bearer ' . $auth_json["token"];
        $login_data = [
            "grantType" => "password",
            "password" => $this->PASSWORD,
            "language" => "id"
        ];
        $this->getCURL($this->URL_AUTH, json_encode($login_data), $login_authorization);
    }

    function crawlListData()
    {
        $data = $this->getDataJson($this->getCURL($this->URL_TARGET_CRAWL));
        $data = $this->json->decode($data);
        $data = $data->states->items->collections;
        $data = get_object_vars($data);
        $data = reset($data);
        $list_data = [];
        foreach ($data as $item) {
            if (!$this->spreedSheet->isDataExist($item, "olx", "A:A")) {
                array_push($list_data, $this->crawlDetail($item));
            }
        }
        /*for ($i = 0; $i < 1; $i++) {
            array_push($list_data, crawlDetail($data[$i]));
        }*/

        if (count($list_data) > 0) $this->spreedSheet->sendData($list_data);
    }

    function crawlDetail($item)
    {
        $url = "$this->BASE_URL/item/$item";
        $html = $this->getDataJson($this->getCURL($url));
        $html = $this->json->decode($html);
        $elements = $html->states->items->elements->{$item};
        $idUser = $elements->user_id;

        //Data User
        $htmlUser = $html->states->users->elements->{$idUser};
        $name = $htmlUser->name;
        $profile = "$this->BASE_URL/profile/$idUser";
        $locationUser = $htmlUser->locations[0];
        $locationUser = "https://maps.google.com/?q=$locationUser->lat,$locationUser->lon";
        $badges = [];
        foreach ($htmlUser->badges as $item) {
            if ($item->status == 1) array_push($badges, $item->name);
        };
        $badges = implode(", ", $badges);

        //Data items
        $id = $elements->id;
        $url = $html->props->location->hostname . $html->props->location->pathname;
        $created_at = $elements->created_at;
        $title = $elements->title;
        $description = $elements->description;
        $location = $elements->locations[0];
        $location = "https://maps.google.com/?q=$location->lat,$location->lon";
        $price = $elements->price->value->display;
        $images = [];
        $thumbnail = "";
        for ($i = 0; $i < count($elements->images) && $i < 3; $i++) {
            if ($i == 0) $thumbnail = $elements->images[$i]->url;
            array_push($images, "_" . $elements->images[$i]->url . "_");
        }
        $images = implode(",\n", $images);

        $parameters = [];
        foreach ($elements->parameters as $item) {
            array_push($parameters, "$item->key_name : $item->formatted_value");
            if ($item->key_name == "phone") {
                $phone = explode("+", $item->formatted_value)[1];
                array_push($parameters, "wa : _https://wa.me/$phone" . "_");
            }
        };
        $parameters = implode("\n", $parameters);

        return $data = [
            "id" => $id,
            "user" => $name,
            "profile" => $profile,
            "locationUser" => $locationUser,
            "badges" => $badges,
            "title" => $title,
            "created_at" => $created_at,
            "url" => $url,
            "location" => $location,
            "thumbnail" => $thumbnail,
            "image" => $images,
            "price" => $price,
            "parameters" => $parameters,
            "description" => $description
        ];

    }

    function getDataJson($html)
    {
        $html = trim(preg_replace('/[\t\n\r\s]+/', ' ', $html));
        $html = explode('<script type="text/javascript">', $html);
        if (!$html[3]) print_r($html);
        $html = explode("; </script>", $html[3]);
        $html = explode('window.__APP =', $html[0]);
        return $html[1];
    }
}

new CrawlDataOLX();