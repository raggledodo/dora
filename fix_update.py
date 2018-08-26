#!/usr/bin/env python3

''' Fix Update '''

import os
import re

THIS_DIR = os.path.dirname(os.path.realpath(__file__))
REMOVE_LINES = [
    '"//vendor/golang.org/x/net/context:go_default_library"',
    '"//vendor/google.golang.org/grpc:go_default_library"',
	'"@io_bazel_rules_go//proto/wkt:empty_go_proto"'
]
QUERY_PATTERN = r'deps\s*=\s*\[\s*((?:.*\n)+)\s*\],'
REPLACE_CONTENT = '''deps = [%s] + [
		"@com_github_golang_protobuf//ptypes/empty:go_default_library",
	] + GRPC_COMPILE_DEPS,'''
FILE = os.path.join(THIS_DIR, "server/BUILD.bazel")

with open(FILE, 'r') as src:
	content = src.read()
	out = re.search(QUERY_PATTERN, content)
	groups = None
	if out:
		groups = out.groups()
	if groups is not None and len(groups) > 0:
		lines = groups[0].split('\n')
		lines = [s.strip().replace(',', '') for s in lines]
		lines = list(filter(lambda s : s not in REMOVE_LINES and len(s) > 0, lines))
		replacement = REPLACE_CONTENT % ('\n        ' +
										 ',\n        '.join(lines) +
										 ',\n    ')
		content = content.replace(out.group(0), replacement)

lines = [
	'load("@org_pubref_rules_protobuf//go:rules.bzl", "GRPC_COMPILE_DEPS")',
]
for line in lines:
	if line not in content:
		content = line + content

if content:
	with open(FILE, 'w') as dest:
		dest.write(content)
