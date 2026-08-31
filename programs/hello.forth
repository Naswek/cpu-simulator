
: puts
  begin
    dup @ 0= not
  while
    dup @ emit
    char+
  repeat
  drop
;

s" Hello, World!" puts
