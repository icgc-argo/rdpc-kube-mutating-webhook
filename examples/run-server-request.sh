#!/bin/bash
#
# Copyright (c) 2020. Ontario Institute for Cancer Research
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.
#

port=$1
response=$(curl -s -XPOST \
	-H 'Content-Type: application/json' \
	-d "@admission-review.example.json" \
	"http://localhost:${port}/mutate" | jq . )

echo "Response: "
echo "$response"
echo ""
jsonpatch=$(echo "$response" | jq -r .response.patch | base64 -d | jq .)
echo "Decoded jsonPatch Response: "
echo "$jsonpatch"

