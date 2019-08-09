SHELL = /bin/bash

generated_img_dir=img/generated

tikz_files := $(filter-out img-src/_%,$(wildcard img-src/*.tikz))

tikz_pngs := $(addprefix $(generated_img_dir)/,$(addsuffix .png,$(notdir $(basename $(tikz_files)))))

.PHONY: clean docs pngs

clean:
	shopt -s nullglob; \
	rm -f build.tex *.aux *.log *.gz *.pdf *.ps *.dvi *.out *.fls *.fdb_latexmk \
		img-src/*.aux img-src/*.log img-src/*.gz img-src/*.pdf img-src/*.ps \
		img-src/*.dvi img-src/*.out img-src/*.fls img-src/*.fdb_latexmk; \
	rm -rf img/generated

$(tikz_pngs): img/generated/%.png: img-src/%.tikz
	bash scripts/generate-pngs.sh /tmp/png-generate img/generated $< $(basename $(*F))

pngs: $(tikz_pngs)

build.tex: $(tikz_files)
	bash scripts/generate-tex.sh $^

build.pdf: build.tex
	pdflatex -halt-on-error build.tex; \
	pdflatex -halt-on-error build.tex

docs: build.pdf pngs
