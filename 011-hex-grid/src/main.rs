use std::fs::File;
use std::io::Read;

fn dist(a : (i32, i32, i32), b : (i32, i32, i32)) -> i32 {
    ((a.0 - b.0).abs() + (a.1 - b.1).abs() + (a.2 - b.2).abs()) / 2
}

#[test]
fn dist_test() {
    assert_eq!(dist((0,0,0), (3, 0, -3)), 3);
    assert_eq!(dist((0,0,0), (0, 0, 0)),  0);
    assert_eq!(dist((0,0,0), (2, -2, 0)), 2);
    assert_eq!(dist((0,0,0), (-1, -2, 3)), 3);
}

fn traverse(directions : &Vec<&str>) -> (i32, i32, i32) {
    let mut max_dist = 0;

    let pos = directions.iter().fold((0, 0, 0), |pos, &dir| {
        let x = pos.0;
        let y = pos.1;
        let z = pos.2;

        let new_pos = match dir {
            "ne" => (x+1, y,   z-1),
            "n"  => (x,   y+1, z-1),
            "nw" => (x-1, y+1, z),
            "sw" => (x-1, y,   z+1),
            "s"  => (x,   y-1, z+1),
            "se" => (x+1, y-1, z),
            _ => panic!("unknown direction")
        };

        let d = dist((0, 0, 0), new_pos);

        if d > max_dist {
            max_dist = d;
        }

        new_pos
    });

    println!("{}", max_dist);

    pos
}

#[test]
fn traverse_test() {
    assert_eq!(traverse(&vec!("ne", "ne", "ne")), (3, 0, -3));
    assert_eq!(traverse(&vec!("ne", "ne", "sw", "sw")), (0, 0, 0));
    assert_eq!(traverse(&vec!("ne", "ne", "s", "s")), (2, -2, 0));
    assert_eq!(traverse(&vec!("se", "sw", "se", "sw", "sw")), (-1, -2, 3));
}

fn main() {
    let mut content = String::new();
    let mut file = File::open("input.txt").expect("file not found");

    file.read_to_string(&mut content).expect("can't read file");

    let steps = content.trim().split(",").collect();

    let d = dist((0, 0, 0), traverse(&steps));

    println!("{}", d);
}
