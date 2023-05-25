let reduce = fn(arr, init, f) {
    let iter = fn(arr, result) {
        if (len(arr) == 0) {
            return result;
        } else {
            return iter(rest(arr), f(result, first(arr)));
        }
    };

    return iter(arr, init);
};

let sum = fn(arr) {
    return reduce(arr, 0, fn(init, acc) { init + acc }); 
};

sum([1, 2, 3, 4, 5, 6, 7, 8, 9]);
