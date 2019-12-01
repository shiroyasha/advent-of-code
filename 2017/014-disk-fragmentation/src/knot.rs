fn reverse(array : &mut Vec<i32>, position: i32, size: i32) {
    let mut iterations = size / 2;
    let mut start = position;
    let mut end = (position + size - 1) % (array.len() as i32);

    while iterations > 0 {
        let tmp = array[start as usize];

        array[start as usize] = array[end as usize];
        array[end   as usize] = tmp;

        start = (start + 1 + (array.len() as i32)) % (array.len() as i32);
        end   = (end   - 1 + (array.len() as i32)) % (array.len() as i32);

        iterations -= 1;
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

fn dense(array : &Vec<i32>) -> String {
    let mut result = "".to_string();

    for i in 0..16 {
        let mut d = 0;

        for j in 0..16 {
            d = d ^ array[i*16 + j];
        }

        if d < 16 {
            result = format!("{}0{:x}", result, d);
        } else {
            result = format!("{}{:x}", result, d);
        }
    }

    result
}

pub fn hash(input: &str) -> String {
    let salt = [17, 31, 73, 47, 23];

    let mut current_pos = 0;
    let mut skip_size = 0;
    let mut array : Vec<i32> = Vec::new();
    let mut lengths = Vec::new();

    for i in 0..256 {
        array.push(i);
    }

    for c in input.bytes() {
        lengths.push(c as i32);
    }

    for s in salt.iter() {
        lengths.push(*s);
    }

    for _ in 0..64 {
        for l in &lengths {
            // println!("{:?}, {}, {}, {}", array, current_pos, skip_size, l);

            reverse(&mut array, current_pos, *l);

            current_pos = (current_pos + l + skip_size) % (array.len() as i32);
            skip_size += 1;
        }
    }

    dense(&array)
}
