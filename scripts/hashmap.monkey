let map = {
    "test1": 1234,
    "test2": fn(x) {x},
    1234   : true,
    false  : 0,
};

puts(map["test1"] == 1234)
puts(map["test2"](42) == 42)
puts(map[1234] == true)
puts(map[false] == 0)
