use Crypt::Mode::ECB;
use MIME::Base64;

print "Enter a string: ";
my $user_input = <STDIN>;
chomp $user_input;

my $regex = qr/^brics\+{[A-Za-z]{21}}$/mp;
if ( !($user_input =~ /$regex/g) ) {
    print "Invalid flag\n";
    exit;
}

$flag = substr($user_input, 7, -1);
if (length($flag) != 21) {
    print "Invalid flag\n";
    exit;
}

@prgs = (
'u4X3gXXUuD1R7O6Wx8ElQIitxCKLVJXPzlUQmCGnl0bjdMqfw6SaIYJ4dP5IAgYe7pa6lvzI3zzvJwWhgq5xjgWR8EWlb8pVjmBOfhqD2eQ=',
'FlfV5JzfiojM0ZFBylmwg4itxCKLVJXPzlUQmCGnl0YiG4zI6M4SBU0AEE+DhLEf7pa6lvzI3zzvJwWhgq5xjgWR8EWlb8pVjmBOfhqD2eQ=',
'hYOrwQOoocIhkAdAA3AIYsKFvxqwWy2wkeDAqVmFo0+Eww0opU1s0S1MNkf1aOkOWkllY23a6e8KeceHAENbEvvm7PCgSoffCgaLP/w+tZM=',
'MpNJzKHD0baGtLeexkKUG8KFvxqwWy2wkeDAqVmFo0+XinWhyHDZL+QpcKZnWbhKguJjX0n23Q0i54n6V5dwShg8wnsdQEHJVrCf6EZhw+M=',
'tSWFrOUBZdCtI0CRDYJGZ8/54mdPkEVXGQewaFlZNLQ48LQG9ct6RMJeMXoSzgnhImVkcskYq0jQLqUxAAzqzGq5ns50Hq7cDJtC7nP0PMM=',
'UrqrneCBf05h/0w4Ebq+B8/54mdPkEVXGQewaFlZNLQLjlgmwFIcGQhkztIODWPRImVkcskYq0jQLqUxAAzqzGq5ns50Hq7cDJtC7nP0PMM=',
'qyBXd3FC0Rn+rsX6EbSl+qankHYOvZHzJ7GPZKyPiR0ozBfagrjGPHToefXZvA3Q847xK/9Uiq8HH+6Y0nr+HIU0p2NHASPfdB/+LSiYkl0=',
'NWthaJjlHmVpuVJJT7D4c6ankHYOvZHzJ7GPZKyPiR3YeBw9MEWaFEZcy6eNEGWk847xK/9Uiq8HH+6Y0nr+HIU0p2NHASPfdB/+LSiYkl0=',
'0q8uGc9CN/WNkiadhK/1vgzndiqk65xzDDI4P5ED/5PlUWCZlUAmAY9GHg0a1L7VW2k5JaL/CiKgKAXQqnkK3iMZMOtroXhLuXBxIVqSW/s=',
'92L7X5izT3gx9yaOiMqIRgzndiqk65xzDDI4P5ED/5M5wGjo1eaA7B2Fgjb+h7eCw98APZet3rIjKy8a9IKALcIJcM8aAhJaKz0XmXXagFQ=',
'ELJpYkI+0ngxtQMyaycdC9xUrMCqc5Cz85fwJ9HVBqjECSgp/pTPdwbyRruzBzrFQn2qzkPDmvozcGNlvabiJ7B80YTO7J9ZzvtiYHqgCl8=',
'D65Xua4cDZNgA1Ja2AL4HNxUrMCqc5Cz85fwJ9HVBqinDiVK1TTzgPC7slthzYVDQn2qzkPDmvozcGNlvabiJ7B80YTO7J9ZzvtiYHqgCl8=',
'8wYnO+pldMFZrGoek6ADYNyunhFqaw9ALGl99M8Xm8v114BlzoitOIAmCPVWAZKIjgYQXIMazrOSFzKqRqaADg//tbqOzF8+NuWWB1R9wuE=',
'LWdupdd1JzfXGvO43RFYeNyunhFqaw9ALGl99M8Xm8vIZFpSYVTypbF1yHbd3BWPjgYQXIMazrOSFzKqRqaADg//tbqOzF8+NuWWB1R9wuE=',
'vKZrXi95dZwRLLPl4ufHlEhO7gAsqkpI9MlC5Gwt7rE6n/t6LZeAW0H9GQZ1nKe4SBGHm7ZeH1264NYdS2p8qlgoV81i6zLH/BEcGD/ggv4=',
'i+5FnPfQ1jHrYCzjVxFP/UhO7gAsqkpI9MlC5Gwt7rG1OTIDUGufSkEo3aqAZ2M5kZNAoxVK5xNuBfKOsBNHXOp/LLV86wKYdXbf4otpxXs=',
'eGe9m+MtXgXFScPgXcpbcOAa2H6t7LgpJr5kMT+5wLmNIcoBGfp2pVJMtSY6H/ftVFHXIYtNlpJMuTy7sUTvVusO1xCaXD8fSeoxcmtmXUo=',
'MrmXGbPv179QuPxvFti1v+Aa2H6t7LgpJr5kMT+5wLkxe5vzTqIZxQBybxDaYdS1VFHXIYtNlpJMuTy7sUTvVusO1xCaXD8fSeoxcmtmXUo=',
'6c2/SP1AWMCeAaBjg7qmhTAN1uoKm/mz/8KhKFCGMu0iFtJ1Plh5ZPVGjKC7A9fl6bAmGi5hvBbO/IajKN0XmsQRXbltx19eog+xGFwfO38=',
'KlkX7AdGl5uCfn/HXWRKPjAN1uoKm/mz/8KhKFCGMu3/eDIJ1oGvJQJpWA4AjNMfjhL1bqBXzmPjm6U19HiC2MKpDN7PQkPr+CHjq6AZcXk=',
'PET/FAzMn+uspd6g38Af+2v2ZP3mQH8UoCdCvpif2RgBBhdleIwE7gox88vNi6K23YtQpQCW3rMf+Hmh8yVYd0cPB67O34aC07rlE42EX6A=',
);

my $xored = 0;
for my $i (0 .. $#prgs)
{
    my $m = Crypt::Mode::ECB->new('AES', 0);
    my $key = substr($flag, $i / 2, 1) x 16;
    my $decoded = decode_base64($prgs[$i]);
    my $plaintext = $m->decrypt($decoded, $key);
    my $prog = sprintf($plaintext, substr($flag, $i, 1));
    my $res =  eval($prog);
    $xored ^= $res;
}
if ($xored == 55871) {
    print "Correct!\n";
} else {
    print "Invalid flag!\n";
}
