package backup_test

import (
	"io/ioutil"
	"os"

	"github.com/greenplum-db/gpbackup/backup"
	"github.com/greenplum-db/gpbackup/testutils"
	"github.com/greenplum-db/gpbackup/utils"

	. "github.com/onsi/ginkgo"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var _ = Describe("backup/data tests", func() {
	var connection *utils.DBConn
	var mock sqlmock.Sqlmock
	BeforeEach(func() {
		connection, mock = testutils.CreateAndConnectMockDB()
		testutils.SetupTestLogger()
	})
	Describe("CopyTableOut", func() {
		It("will dump a table to its own file", func() {
			testTable := backup.Relation{2345, 3456, "public", "foo", []string{}}
			execStr := "COPY public.foo TO '<SEG_DATA_DIR>/backups/20170101/20170101010101/gpbackup_<SEGID>_20170101010101_3456' WITH CSV DELIMITER ',' ON SEGMENT;"
			mock.ExpectExec(execStr).WillReturnResult(sqlmock.NewResult(10, 0))
			filename := "<SEG_DATA_DIR>/backups/20170101/20170101010101/gpbackup_<SEGID>_20170101010101_3456"
			backup.CopyTableOut(connection, testTable, filename)
		})
	})
	Describe("WriteTableMapFile", func() {
		testutils.SetDefaultSegmentConfiguration()
		tableOne := backup.Relation{0, 1234, "public", "foo", []string{}}
		tableTwo := backup.Relation{0, 2345, "public", "foo|bar", []string{}}

		It("writes a map file containing one table", func() {
			filePath := ""
			r, w, _ := os.Pipe()
			utils.System.OpenFile = func(name string, flag int, perm os.FileMode) (*os.File, error) { filePath = name; return w, nil }
			defer func() { utils.System.OpenFile = os.OpenFile }()
			tables := []backup.Relation{tableOne}
			backup.WriteTableMapFile(tables)
			w.Close()
			output, _ := ioutil.ReadAll(r)
			testutils.ExpectRegex(string(output), `public.foo: 1234
`)
		})
		It("writes a map file containing multiple tables", func() {
			filePath := ""
			r, w, _ := os.Pipe()
			utils.System.OpenFile = func(name string, flag int, perm os.FileMode) (*os.File, error) { filePath = name; return w, nil }
			defer func() { utils.System.OpenFile = os.OpenFile }()
			tables := []backup.Relation{tableOne, tableTwo}
			backup.WriteTableMapFile(tables)
			w.Close()
			output, _ := ioutil.ReadAll(r)
			testutils.ExpectRegex(string(output), `public.foo: 1234
public."foo|bar": 2345`)
		})
	})
})
