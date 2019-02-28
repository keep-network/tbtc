SHELL = /bin/bash 

generated_img_dir=img/generated

tikz_files := $(wildcard img-src/*.tikz)

tikz_pngs := $(foreach file, $(tikz_files), $(generated_img_dir)/$(notdir $(basename $(file))).png)

.PHONY: clean docs pngs

clean:
	shopt -s nullglob; \
	rm -f *.aux *.log *.gz *.pdf *.ps *.dvi *.out *.fls *.fdb_latexmk \
		img-src/*.aux img-src/*.log img-src/*.gz img-src/*.pdf img-src/*.ps \
		img-src/*.dvi img-src/*.out img-src/*.fls img-src/*.fdb_latexmk; \
	rm -rf img/generated

relay-states.pdf: relay-states.tex
	pdflatex -halt-on-error relay-states.tex; \
	pdflatex -halt-on-error relay-states.tex

$(tikz_pngs): img/generated/%.png: img-src/%.tikz
	bash scripts/generate-pngs.sh /tmp/png-generate img/generated $< $(basename $(*F))

pngs: $(tikz_pngs)

docs: relay-states.pdf $(tikz_pngs)