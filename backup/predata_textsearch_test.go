package backup_test

import (
	"github.com/greenplum-db/gpbackup/backup"
	"github.com/greenplum-db/gpbackup/testutils"
	"github.com/greenplum-db/gpbackup/utils"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("backup/predata_textsearch tests", func() {
	var toc *utils.TOC
	var backupfile *utils.FileWithByteCount

	BeforeEach(func() {
		toc = &utils.TOC{}
		backupfile = utils.NewFileWithByteCount(buffer)
	})
	Describe("PrintCreateTextSearchParserStatements", func() {
		It("prints a basic text search parser", func() {
			parsers := []backup.TextSearchParser{{Oid: 0, Schema: "public", Name: "testparser", StartFunc: "start_func", TokenFunc: "token_func", EndFunc: "end_func", LexTypesFunc: "lextypes_func"}}
			backup.PrintCreateTextSearchParserStatements(backupfile, toc, parsers, backup.MetadataMap{})
			testutils.ExpectEntry(toc.PredataEntries, 0, "public", "testparser", "TEXT SEARCH PARSER")
			testutils.AssertBufferContents(toc.PredataEntries, buffer, `CREATE TEXT SEARCH PARSER public.testparser (
	START = start_func,
	GETTOKEN = token_func,
	END = end_func,
	LEXTYPES = lextypes_func
);`)
		})
		It("prints a text search parser with a headline and comment", func() {
			parsers := []backup.TextSearchParser{{Oid: 1, Schema: "public", Name: "testparser", StartFunc: "start_func", TokenFunc: "token_func", EndFunc: "end_func", LexTypesFunc: "lextypes_func", HeadlineFunc: "headline_func"}}
			metadataMap := testutils.DefaultMetadataMap("TEXT SEARCH PARSER", false, false, true)
			backup.PrintCreateTextSearchParserStatements(backupfile, toc, parsers, metadataMap)
			testutils.AssertBufferContents(toc.PredataEntries, buffer, `CREATE TEXT SEARCH PARSER public.testparser (
	START = start_func,
	GETTOKEN = token_func,
	END = end_func,
	LEXTYPES = lextypes_func,
	HEADLINE = headline_func
);

COMMENT ON TEXT SEARCH PARSER public.testparser IS 'This is a text search parser comment.';`)
		})
	})
	Describe("PrintCreateTextSearchTemplateStatements", func() {
		It("prints a basic text search template with comment", func() {
			templates := []backup.TextSearchTemplate{{Oid: 1, Schema: "public", Name: "testtemplate", InitFunc: "dsimple_init", LexizeFunc: "dsimple_lexize"}}
			metadataMap := testutils.DefaultMetadataMap("TEXT SEARCH TEMPLATE", false, false, true)
			backup.PrintCreateTextSearchTemplateStatements(backupfile, toc, templates, metadataMap)
			testutils.ExpectEntry(toc.PredataEntries, 0, "public", "testtemplate", "TEXT SEARCH TEMPLATE")
			testutils.AssertBufferContents(toc.PredataEntries, buffer, `CREATE TEXT SEARCH TEMPLATE public.testtemplate (
	INIT = dsimple_init,
	LEXIZE = dsimple_lexize
);

COMMENT ON TEXT SEARCH TEMPLATE public.testtemplate IS 'This is a text search template comment.';`)
		})
	})
	Describe("PrintCreateTextSearchDictionaryStatements", func() {
		It("prints a basic text search dictionary with comment", func() {
			dictionaries := []backup.TextSearchDictionary{{Oid: 1, Schema: "public", Name: "testdictionary", Template: "testschema.snowball", InitOption: "language = 'russian', stopwords = 'russian'"}}
			metadataMap := testutils.DefaultMetadataMap("TEXT SEARCH DICTIONARY", false, true, true)
			backup.PrintCreateTextSearchDictionaryStatements(backupfile, toc, dictionaries, metadataMap)
			testutils.ExpectEntry(toc.PredataEntries, 0, "public", "testdictionary", "TEXT SEARCH DICTIONARY")
			testutils.AssertBufferContents(toc.PredataEntries, buffer, `CREATE TEXT SEARCH DICTIONARY public.testdictionary (
	TEMPLATE = testschema.snowball,
	language = 'russian', stopwords = 'russian'
);

COMMENT ON TEXT SEARCH DICTIONARY public.testdictionary IS 'This is a text search dictionary comment.';


ALTER TEXT SEARCH DICTIONARY public.testdictionary OWNER TO testrole;`)
		})
	})
	Describe("PrintCreateTextSearchConfigurationStatements", func() {
		It("prints a basic text search configuration", func() {
			configurations := []backup.TextSearchConfiguration{{Oid: 0, Schema: "public", Name: "testconfiguration", Parser: `pg_catalog."default"`, TokenToDicts: map[string][]string{}}}
			backup.PrintCreateTextSearchConfigurationStatements(backupfile, toc, configurations, backup.MetadataMap{})
			testutils.AssertBufferContents(toc.PredataEntries, buffer, `CREATE TEXT SEARCH CONFIGURATION public.testconfiguration (
	PARSER = pg_catalog."default"
);`)
		})
		It("prints a text search configuration with multiple mappings and comment", func() {
			tokenToDicts := map[string][]string{"int": {"simple", "english_stem"}, "asciiword": {"english_stem"}}
			configurations := []backup.TextSearchConfiguration{{Oid: 1, Schema: "public", Name: "testconfiguration", Parser: `pg_catalog."default"`, TokenToDicts: tokenToDicts}}
			metadataMap := testutils.DefaultMetadataMap("TEXT SEARCH CONFIGURATION", false, true, true)
			backup.PrintCreateTextSearchConfigurationStatements(backupfile, toc, configurations, metadataMap)
			testutils.ExpectEntry(toc.PredataEntries, 0, "public", "testconfiguration", "TEXT SEARCH CONFIGURATION")
			testutils.AssertBufferContents(toc.PredataEntries, buffer, `CREATE TEXT SEARCH CONFIGURATION public.testconfiguration (
	PARSER = pg_catalog."default"
);

ALTER TEXT SEARCH CONFIGURATION public.testconfiguration
	ADD MAPPING FOR "asciiword" WITH english_stem;

ALTER TEXT SEARCH CONFIGURATION public.testconfiguration
	ADD MAPPING FOR "int" WITH simple, english_stem;

COMMENT ON TEXT SEARCH CONFIGURATION public.testconfiguration IS 'This is a text search configuration comment.';


ALTER TEXT SEARCH CONFIGURATION public.testconfiguration OWNER TO testrole;`)
		})
	})
})