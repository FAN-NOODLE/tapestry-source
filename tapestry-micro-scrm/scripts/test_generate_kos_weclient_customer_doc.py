#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import tempfile
import unittest
import xml.etree.ElementTree as ET
from pathlib import Path
from zipfile import ZipFile

from generate_kos_weclient_customer_doc import build_document, validate_document


class CustomerDocumentGeneratorTest(unittest.TestCase):
    def test_builds_valid_five_page_customer_document(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            output = Path(directory) / "客户技术文档.docx"

            build_document(output)
            validate_document(output)

            with ZipFile(output) as archive:
                names = set(archive.namelist())
                document_xml = archive.read("word/document.xml")
                text = document_xml.decode("utf-8")

            self.assertTrue(
                {
                    "[Content_Types].xml",
                    "_rels/.rels",
                    "word/document.xml",
                    "word/styles.xml",
                    "word/settings.xml",
                }.issubset(names)
            )
            ET.fromstring(document_xml)
            self.assertEqual(text.count('w:type="page"'), 4)
            self.assertIn("KOS 与 WeClient 素材数据同步技术说明", text)
            self.assertIn("X-Coach-Mark-Token", text)
            self.assertIn("is_valid = N", text)
            self.assertNotIn("REDACTED_FOR_SOURCE_SHARING", text)
            self.assertNotIn("source_link_kos", text)


if __name__ == "__main__":
    unittest.main()
