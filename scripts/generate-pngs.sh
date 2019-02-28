tmp_dir=$1
img_dir=$(pwd)/$2
tikz_source=$(pwd)/$3
target_basename=$4

mkdir -p $tmp_dir
cd $tmp_dir
echo "\documentclass{standalone}" > $target_basename.tex
echo "\usepackage{tikz}" >> $target_basename.tex
echo "\usepackage{hyperref}" >> $target_basename.tex
echo "\usepackage[utf8]{inputenc}" >> $target_basename.tex
echo "\usetikzlibrary{positioning}" >> $target_basename.tex
echo "\usetikzlibrary{arrows.meta}" >> $target_basename.tex
echo "\usetikzlibrary{shapes.symbols}" >> $target_basename.tex
echo "\usetikzlibrary{calc}" >> $target_basename.tex
echo "\tikzset{every node/.style={above},start state/.style={draw,circle,text width=0},state/.style={draw,circle,align=flush center,text width=2cm},decision/.style={draw,rectangle,align=flush center}, thread/.style={draw,signal,signal to=east,fill=white},nested state/.style={draw,circle,double,align=flush center}, nested decision/.style={draw,rectangle,double,align=flush center}, chain state/.style={draw,circle,dashed,align=flush center,text width=2cm}, chain decision/.style={draw,rectangle,dashed,align=flush center}, chain transition/.style={draw,dashed},nested chain decision/.style={draw,rectangle,double,dashed,align=flush center},>=Stealth }" >> $target_basename.tex
echo "\begin{document}" >> $target_basename.tex
echo "" >> $target_basename.tex
cat $tikz_source >> $target_basename.tex
echo "\end{document}" >> $target_basename.tex

pdflatex -halt-on-error $target_basename.tex $target_basename.pdf

mkdir -p $img_dir
echo "Generating $img_dir/$target_basename.png from $target_basename.pdf..."
pdftoppm $target_basename.pdf $img_dir/$target_basename -png -f 1 -singlefile -rx 300 -ry 300
echo "Finished with exit status $?."
