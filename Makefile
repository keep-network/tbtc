SHELL = /bin/bash

generated_img_dir=img/generated

tikz_files := $(filter-out docs/img-src/_%,$(wildcard docs/img-src/*.tikz))

tikz_pngs := $(foreach file, $(tikz_files), $(generated_img_dir)/$(notdir $(basename $(file))).png)

.PHONY: clean docs pngs

clean:
	shopt -s nullglob; \
	rm -f build.tex *.aux *.log *.gz *.pdf *.ps *.dvi *.out *.fls *.fdb_latexmk \
		docs/img-src/*.aux docs/img-src/*.log docs/img-src/*.gz docs/img-src/*.pdf docs/img-src/*.ps \
		docs/img-src/*.dvi docs/img-src/*.out docs/img-src/*.fls docs/img-src/*.fdb_latexmk; \
	rm -rf img/generated

$(tikz_pngs): img/generated/%.png: docs/img-src/%.tikz
	bash scripts/generate-pngs.sh /tmp/png-generate img/generated $< $(basename $(*F))

pngs: $(tikz_pngs)

build.tex: $(tikz_files)
	bash scripts/generate-tex.sh $^

build.pdf: build.tex
	pdflatex -halt-on-error build.tex; \
	pdflatex -halt-on-error build.tex

docs: build.pdf pngs
