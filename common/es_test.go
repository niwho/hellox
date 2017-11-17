package common

import (
	"os"
	"testing"

	elastic "gopkg.in/olivere/elastic.v3"
	"pusic/push/kafka2es/logs"
)

func TestMain(m *testing.M) {

	// InitEsClient([]string{"http://10.116.27.15:9200", "http://10.116.27.33:9200"}, "test_wt", "tt", 5, 0, 10)
	InitUdpClient("9067aadfbee458e0d2e1f9876e898882aab5f5876f617634e7e84ec9be5bded0")
	InitEsClient([]string{"http://119.81.218.90:5858"}, "pusic_push_test", "gopush", 5, 0, 10)
	os.Exit(m.Run())
}

func Test_Insert(t *testing.T) {
	//AnchorOnlinePush2es(AnchorOnShow{"11111", "", "success", "http://10.116.27.33:9200", 123, "niwho", "", 666, 0})
}

func Test_Search(t *testing.T) {
	q := elastic.NewRawStringQuery(`{ "from": 0, "size": 0, "_source": { "includes": [ "COUNT" ], "excludes": [] }, "aggregations": { "cnt": { "value_count": { "field": "ts" } } } }`)
	rt, err := EsClient.Search("pusic_push_test_2017_10_24").Query(q).Do()
	logs.Log(logs.F{"rt": rt, "err": err}).Info("query test")

}
