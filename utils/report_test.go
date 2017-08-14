package utils_test

import (
	"strings"

	"github.com/greenplum-db/gpbackup/testutils"
	"github.com/greenplum-db/gpbackup/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("utils/report tests", func() {
	Describe("ParseErrorMessage", func() {
		It("Parses a CRITICAL error message and returns error code 1", func() {
			var err interface{}
			err = "testProgram:testUser:testHost:000000-[CRITICAL]:-Error Message"
			errMsg, exitCode := utils.ParseErrorMessage(err)
			Expect(errMsg).To(Equal("Error Message"))
			Expect(exitCode).To(Equal(1))
		})
		It("Returns error code 0 for an empty error message", func() {
			errMsg, exitCode := utils.ParseErrorMessage(nil)
			Expect(errMsg).To(Equal(""))
			Expect(exitCode).To(Equal(0))
		})
	})
	Describe("WriteReportFile", func() {
		timestamp := "20170101010101"
		backupReport := utils.Report{DatabaseName: "testdb", DatabaseVersion: "5.0.0 build test", BackupVersion: "0.1.0", BackupType: "Unfiltered Full Backup"}
		objectCounts := map[string]int{"tables": 42, "sequences": 1, "types": 1000}

		It("writes a report for a successful backup", func() {
			utils.WriteReportFile(connection, buffer, timestamp, backupReport, objectCounts, "42 MB", "")
			Expect(buffer).To(gbytes.Say(`Greenplum Database Backup Report

Timestamp Key: 20170101010101
GPDB Version: 5\.0\.0 build test
gpbackup Version: 0\.1\.0

Database Name: testdb
Command Line: .*
Backup Type: Unfiltered Full Backup
Backup Status: Success

Database Size: 42 MB
Count of Database Objects in Backup:
sequences                	1
tables                   	42
types                    	1000`))
		})
		It("writes a report for a failed backup", func() {
			utils.WriteReportFile(connection, buffer, timestamp, backupReport, objectCounts, "42 MB", "Cannot access /tmp/backups: Permission denied")
			Expect(buffer).To(gbytes.Say(`Greenplum Database Backup Report

Timestamp Key: 20170101010101
GPDB Version: 5\.0\.0 build test
gpbackup Version: 0\.1\.0

Database Name: testdb
Command Line: .*
Backup Type: Unfiltered Full Backup
Backup Status: Failure
Backup Error: Cannot access /tmp/backups: Permission denied

Database Size: 42 MB
Count of Database Objects in Backup:
sequences                	1
tables                   	42
types                    	1000`))
		})
	})
	Describe("ReadReportFile", func() {
		It("can read a report file for a successful backup", func() {
			reportFileContents := `Greenplum Database Backup Report

Timestamp Key: 20170101010101
GPDB Version: 5.0.0 build test
gpbackup Version: 0.1.0

Database Name: testdb
Command Line: gpbackup --dbname testdb
Backup Type: Unfiltered Full Backup
Backup Status: Success

Database Size: 42 MB
Count of Database Objects in Backup:
sequences                   1
tables                      42
types                       1000
`
			reportReader := strings.NewReader(reportFileContents)
			backupReport := utils.ReadReportFile(reportReader)
			expectedReport := utils.Report{DatabaseName: "testdb", DatabaseVersion: "5.0.0 build test", BackupVersion: "0.1.0", BackupType: "Unfiltered Full Backup"}
			testutils.ExpectStructsToMatch(&expectedReport, &backupReport)
		})
		It("can read a report file for a failed backup", func() {
			reportFileContents := `Greenplum Database Backup Report

Timestamp Key: 20170101010101
GPDB Version: 5.0.0 build test
gpbackup Version: 0.1.0

Database Name: testdb
Command Line: gpbackup --dbname testdb
Backup Type: Unfiltered Full Backup
Backup Status: Failure
Backup Error: Cannot access /tmp/backups: Permission denied

Database Size: 42 MB
Count of Database Objects in Backup:
sequences                   1
tables                      42
types                       1000
`
			reportReader := strings.NewReader(reportFileContents)
			backupReport := utils.ReadReportFile(reportReader)
			expectedReport := utils.Report{DatabaseName: "testdb", DatabaseVersion: "5.0.0 build test", BackupVersion: "0.1.0", BackupType: "Unfiltered Full Backup"}
			testutils.ExpectStructsToMatch(&expectedReport, &backupReport)
		})
	})
})
