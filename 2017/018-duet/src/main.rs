use std::fs::File;
use std::io::Read;

mod program;

use program::{Program,State};

fn transfer(p1 : &mut Program, p2 : &mut Program) {
    for v in &p1.sent {
        p2.incomming.insert(0, v.clone());
    }

    p1.sent = Vec::new();
}

fn is_finished(p1 : &Program, p2 : &Program) -> bool {
    println!("STATES: {:?}, {:?}, {:?}, {:?}", p1.state, p2.state, p1.incomming, p2.incomming);

    if p1.state == State::Finished && p2.state == State::Finished {
        return true;
    }

    // deadlock
    if p1.state == State::Locked && p1.incomming.len() == 0 && p2.state == State::Locked && p2.incomming.len() == 0 {
        return true;
    }

    false
}

fn execute(source : &str) -> i64 {
    let mut p1 = Program::new(0, source);
    let mut p2 = Program::new(1, source);

    while !is_finished(&p1, &p2) {
        p1.run();

        println!("BEFORE TRANSFER: {:?}", p1.sent);

        transfer(&mut p1, &mut p2);

        println!("AFTER TRANSFER: {:?}", p2.incomming);

        p2.run();
        transfer(&mut p2, &mut p1);
    }

    p2.sent_count
}

fn main() {
    let mut content = String::new();
    let mut file = File::open("input.txt").expect("file not found");

    file.read_to_string(&mut content).expect("can't read file");

    println!("{}", execute(&content));
}

#[test]
fn transfer_test() {
    let mut p1 = Program::new(0, "");
    let mut p2 = Program::new(1, "");

    p1.sent.push(1);
    p1.sent.push(2);
    p1.sent.push(3);

    p2.incomming.push(4);
    p2.incomming.push(5);

    transfer(&mut p1, &mut p2);

    assert_eq!(p1.sent, vec![]);
    assert_eq!(p2.incomming, vec![3, 2, 1, 4, 5]);
}

#[test]
fn execute_test() {
    let source = "snd 1
snd 2
snd p
rcv a
rcv b
rcv c
rcv d
";

   assert_eq!(execute(source), 3);
}
