find docs -name "*.png" | xargs docker run --rm -t -v $PWD:/source luizeof/image-optimizer
