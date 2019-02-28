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

tbtc-diagrams.pdf: tbtc-diagrams.tex
	pdflatex -halt-on-error tbtc-diagrams.tex; \
	pdflatex -halt-on-error tbtc-diagrams.tex

$(tikz_pngs): img/generated/%.png: img-src/%.tikz
	bash scripts/generate-pngs.sh /tmp/png-generate img/generated $< $(basename $(*F))

pdfs: tbtc-diagrams.pdf

pngs: $(tikz_pngs)

docs: tbtc-diagrams.pdf $(tikz_pngs)
