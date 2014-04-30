#! /bin/bash

for tutorial in tut*/; do
	cd $tutorial
	for chapter in */; do
		cd $chapter
		testfile="${chapter}_test.go"
		if [ ! -e "$testfile" ]; then
			cat >"$testfile"<<<"package $chapter

import \"testing\"

func TestRun(t *testing.T) {
	Run()
}"
		fi
		cd ..
	done
	cd ..
done

