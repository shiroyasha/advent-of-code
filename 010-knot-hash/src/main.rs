
fn reverse(array : &mut Vec<i32>, position: i32, size: i32) {
    let mut iterations = size / 2;
    let mut start = position;
    let mut end = (position + size - 1) % (array.len() as i32);

    while iterations > 0 {
        // println!("{}: {:?}, {}, {}", iterations, array, start, end);

        let tmp = array[start as usize];

        array[start as usize] = array[end as usize];
        array[end   as usize] = tmp;

        start = (start + 1 + (array.len() as i32)) % (array.len() as i32);
        end   = (end   - 1 + (array.len() as i32)) % (array.len() as i32);

        iterations -= 1;

        // println!("{}: {:?}, {}, {}", iterations +1, array, start, end);
    }
}

#[test]
fn reverse_test() {
    let mut a = vec![0, 1, 2, 3, 4];

    reverse(&mut a, 0, 3);

    assert_eq!(a, vec![2, 1, 0, 3, 4]);

    let mut b = vec![0, 1, 2, 3, 4];

    reverse(&mut b, 3, 4);

    assert_eq!(b, vec![4, 3, 2, 1, 0]);
}

fn hash(array_size : i32, lengths: &Vec<i32>) -> i32 {
    let mut current_pos = 0;
    let mut skip_size = 0;

    let mut array : Vec<i32> = Vec::new();

    // initialize the list
    for i in 0..array_size {
        array.push(i);
    }

    for l in lengths {
        println!("{:?}, {}, {}, {}", array, current_pos, skip_size, l);

        reverse(&mut array, current_pos, *l);

        current_pos = (current_pos + l + skip_size) % (array.len() as i32);
        skip_size += 1;
    }

    println!("{:?}, {}, {}", array, current_pos, skip_size);

    array[0] * array[1]
}

#[test]
fn hash_test() {
    let lengths = vec![3, 4, 1, 5];

    assert_eq!(hash(5, &lengths), 12);
}

fn main() {
    let lengths = vec![192,69,168,160,78,1,166,28,0,83,198,2,254,255,41,12];

    println!("{}", hash(256, &lengths));
}
