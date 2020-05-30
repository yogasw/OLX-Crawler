<?php

use PhpAmqpLib\Connection\AMQPStreamConnection;
use PhpAmqpLib\Message\AMQPMessage;

require __DIR__ . "/vendor/autoload.php";

class SpreedSheet
{
    public $service;
    public $spreadsheetId = "1a99s1jK3T-wBZI1rlpLsECdi5WFnpq8h3U6uUP8oNIs";

    function __construct()
    {
        $client = new \Google_Client();
        $client->setApplicationName("Google Sheets and PHP");
        $client->setScopes([\Google_Service_Sheets::SPREADSHEETS]);
        $client->setAccessType('offline');
        $client->setAuthConfig(__DIR__ . '/elated-lotus-161202-10ba0a49a397.json');
        $this->service = new Google_Service_Sheets($client);
    }

    private function update($values)
    {
        /*$response = $service->spreadsheets_values->get($spreadsheetId, $range);
        $values = $response->getValues();
        $values = [
            ["A1", "B1"]
        ];*/

        $range = "olx!A2:F2";
        $params = [
            'valueInputOption' => 'RAW'
        ];
        $body = new Google_Service_Sheets_ValueRange([
            'values' => $values
        ]);

        $this->service->spreadsheets_values->update(
            $this->spreadsheetId,
            $range,
            $body,
            $params
        );
    }

    private function insert($item)
    {
        $range = "olx!A2:I2";
        $values = [
            $item["id"],
            $item["title"],
            $item["created_at"],
            $item["url"],
            $item["location"],
            $item["image"],
            $item["price"],
            $item["parameters"],
            $item["description"]
        ];
        $body = new Google_Service_Sheets_ValueRange([
            'values' => [$values]
        ]);

        $insert = [
            'valueInputOption' => 'USER_ENTERED'
        ];
        $this->service->spreadsheets_values->append(
            $this->spreadsheetId,
            $range,
            $body,
            $insert
        );
    }

    private function isDataExist($searchValue, $searchSheet, $searchRange)
    {
        $params = [
            'valueInputOption' => 'USER_ENTERED',
            'includeValuesInResponse' => true,
            'responseValueRenderOption' => true,
        ];
        $range = "HIDE_SHEET!A1:B1";
        $values = [
            ['=MATCH(' . $searchValue . ', ' . $searchSheet . '!' . $searchRange . ', 0)']
        ];
        $body = new Google_Service_Sheets_ValueRange([
            'values' => $values
        ]);
        $result = $this->service->spreadsheets_values->update(
            $this->spreadsheetId,
            $range,
            $body,
            $params
        );

        try {
            $value = intval($result->updatedData->values[0][0]);
            if ($value > 0) {
                return true;
            } else {
                return false;
            }
        } catch (\mysql_xdevapi\Exception $e) {
            return false;
        }

    }

    function sendData($data)
    {
        $connection = new AMQPStreamConnection('moose.rmq.cloudamqp.com', 5672, getenv("RABBITMQ_DEFAULT_USER"), getenv("RABBITMQ_DEFAULT_PASS"), getenv("RABBITMQ_DEFAULT_VHOST"));
        $channel = $connection->channel();
        if (!$connection->isConnected()) {
            throw new Exception('Connection RabbitMQ Failed. \n');
        } else {
            print_r("Connection RabbitMQ Success \n");
        }
        $channel->queue_declare(getenv("RABBITMQ_DEFAULT_QUEUE"), false, false, false, false);

        foreach ($data as $item) {
            if (!$this->isDataExist($item["id"], "olx", "A:A")) {
                $this->insert($item);
                $message = "*{$item["title"]}*\n\n"
                    . "*Ad description*\n"
                    . "*Created At :*\n{$item["created_at"]}\n\n"
                    . "*Location :*\n{$item["location"]}\n\n"
                    . "*Image :*\n{$item["image"]}\n\n"
                    . "*Price :* {$item["price"]}\n\n"
                    . "*Parameters :* \n{$item["parameters"]}\n\n"
                    . "*Description :* \n{$item["description"]}\n\n"
                    . "*Url :* \n{$item["url"]}_\n\n"
                    . "*User description :* \n"
                    . "*Name :* {$item["user"]}\n\n"
                    . "*Profile :*\n{$item["profile"]}\n\n"
                    . "*Location :*\n{$item["locationUser"]}\n\n"
                    . "*Badges :* {$item["badges"]}\n\n";
                $message = [
                    "Target" => "6282329949292-1590306644@g.us",
                    "Message" => $message
                ];
                $msg = new AMQPMessage(json_encode($message));
                $channel->basic_publish($msg, '', 'wa-text');
                print_r("send message id :" . $item['id'] . "\n");
            }
        }

        $channel->close();
        $connection->close();
        print_r("Connection RabbitMQ Closed \n");
    }
}
