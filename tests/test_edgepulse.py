import csv
import io
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

import sys
sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from edgepulse import (
    Endpoint,
    format_table,
    load_targets,
    parse_targets,
    write_csv,
)


class EdgePulseTests(unittest.TestCase):
    def test_parse_targets_accepts_addresses_and_cidrs(self):
        targets = parse_targets("1.1.1.1, 2001:db8::/126")
        self.assertEqual([str(t) for t in targets], ["1.1.1.1", "2001:db8::/126"])

    def test_load_targets_ignores_comments_and_blank_lines(self):
        with tempfile.TemporaryDirectory() as folder:
            path = Path(folder) / "targets.txt"
            path.write_text("# sample\n1.1.1.1\n\n2001:db8::1\n", encoding="utf-8")
            self.assertEqual([str(t) for t in load_targets(path)], ["1.1.1.1", "2001:db8::1"])

    def test_csv_export_has_english_headers(self):
        row = Endpoint("1.1.1.1", sent=3, received=3, loss=0.0, latency_ms=12.5, speed_mbps=4.2, region="LAX")
        output = io.StringIO()
        write_csv([row], output)
        output.seek(0)
        records = list(csv.reader(output))
        self.assertEqual(records[0], ["Address", "Sent", "Received", "Loss %", "Latency (ms)", "Speed (MB/s)", "Region"])
        self.assertEqual(records[1][0], "1.1.1.1")

    def test_table_is_english(self):
        row = Endpoint("1.1.1.1", sent=1, received=1, loss=0.0, latency_ms=8.0, speed_mbps=1.0, region="N/A")
        self.assertIn("Address", format_table([row]))
        self.assertNotRegex(format_table([row]), r"[\u4e00-\u9fff]")


if __name__ == "__main__":
    unittest.main()
