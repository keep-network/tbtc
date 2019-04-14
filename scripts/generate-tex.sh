target_basename=build
tikz_include_base=$(dirname $(ls $@ | head -n 1))

metadata={}
for tikz_source in $@; do
    front_matter=$(
        cat $tikz_source |
        ruby -e "puts STDIN.take_while{|ln| ln.start_with?('%')}.join" |
        grep "% *![a-z]* =" |
        sed -e "s/% *!//"
    )
    new_metadata=$(
        echo "$front_matter" |
        ruby -e 'puts STDIN.inject({}) { |map, ln| (prop,val) = ln.strip.split(/ *= */); map[prop] = val; map }'
    )
    metadata=$(
        echo "\
            ${new_metadata}.inject(${metadata})  do |final, (key, val)|
                new_val = [final[key],val].compact.join(',')
                final[key] = new_val
                final
            end
        " |
        ruby -e 'puts eval(STDIN.read)'
    )
done

tex_includes=$(
    echo "$metadata" |
    ruby -e 'puts (eval(STDIN.read)["include"] || "").split(",").uniq.join("\n")'
)
tex_package_includes=$(
    echo "$metadata" |
    ruby -e 'puts (eval(STDIN.read)["texpackages"] || "").split(",").uniq.join("\n")' |
    sed -e 's/^\(.*\)$/\\usepackage{\1}/'
)
tikz_library_includes=$(
    echo "$metadata" |
    ruby -e 'puts (eval(STDIN.read)["tikzlibraries"] || "").split(",").uniq.join("\n")' |
    sed -e 's/^\(.*\)$/\\usetikzlibrary{\1}/'
)

echo "\documentclass{article}
\usepackage[utf8]{inputenc}
\usepackage[margin=1in]{geometry}
\usepackage{tikz}
\usepackage[colorlinks=true]{hyperref}
\usepackage{varwidth}

$tex_package_includes
$tikz_library_includes

\begin{document}
$(
    for filename in $tex_includes; do
        cat $tikz_include_base/$filename
    done
)

$(
    for tikz_source in $@; do
        tikz_source_base=$(basename $tikz_source)
        tikz_name=$(echo $tikz_source_base | sed -e "s/-\([a-z]\)/\u&/g")

        echo "\begin{figure}
    \centering
    \input{$tikz_source}

    \caption{\label{fig:$tikz_source_base}$tikz_name.}
\end{figure}
"
    done
)

\end{document}" >> $target_basename.tex
