tmp_dir=$1
img_dir=$(pwd)/$2
tikz_source=$(pwd)/$3
tikz_include_base=$(dirname $tikz_source)
target_basename=$4

mkdir -p $tmp_dir
cd $tmp_dir

front_matter=$(
    cat $tikz_source |
    ruby -e "puts STDIN.take_while{|ln| ln.start_with?('%')}.join" |
    grep "% *![a-z]* =" |
    sed -e "s/% *!//"
)
metadata=$(
    echo "$front_matter" |
    ruby -e 'puts STDIN.inject({}) { |map, ln| (prop,val) = ln.strip.split(/ *= */); map[prop] = val; map }'
)
tex_includes=$(
    echo "$metadata" |
    ruby -e 'puts (eval(STDIN.read)["include"] || "").split(",").join("\n")'
)
tex_package_includes=$(
    echo "$metadata" |
    ruby -e 'puts (eval(STDIN.read)["texpackages"] || "").split(",").join("\n")' |
    sed -e 's/^\(.*\)$/\\usepackage{\1}/'
)
tikz_library_includes=$(
    echo "$metadata" |
    ruby -e 'puts (eval(STDIN.read)["tikzlibraries"] || "").split(",").join("\n")' |
    sed -e 's/^\(.*\)$/\\usetikzlibrary{\1}/'
)

echo "\documentclass{standalone}" > $target_basename.tex
echo "\usepackage{tikz}" >> $target_basename.tex
echo "\usepackage{hyperref}" >> $target_basename.tex
echo "\usepackage[utf8]{inputenc}" >> $target_basename.tex
echo "$tex_package_includes" >> $target_basename.tex
echo "$tikz_library_includes" >> $target_basename.tex
echo "" >> $target_basename.tex
echo "\begin{document}" >> $target_basename.tex
echo "" >> $target_basename.tex
for filename in $tex_includes; do
    cat $tikz_include_base/$filename >> $target_basename.tex
done
echo "" >> $target_basename.tex
cat $tikz_source >> $target_basename.tex
echo "\end{document}" >> $target_basename.tex

pdflatex -halt-on-error $target_basename.tex $target_basename.pdf

mkdir -p $img_dir
echo "Generating $img_dir/$target_basename.png from $target_basename.pdf..."
pdftoppm $target_basename.pdf $img_dir/$target_basename -png -f 1 -singlefile -rx 300 -ry 300

finalstatus=$?
echo "Finished with exit status $finalstatus."
exit $finalstatus
