variable limit
variable ready
variable sum
variable f1
variable f2
variable temp

: handle_input
  key dup 10 = if
    drop
    1 ready !
  else
    48 - limit @ 10 * + limit !
  then
  iret
;

: solve
  0 sum !
  1 f1 !
  2 f2 !
  begin
    f2 @ limit @ <
  while
    \ check if f2 is even
    f2 @ 2 mod 0 = if
      sum @ f2 @ + sum !
    then
    \ temp = f1 + f2
    f1 @ f2 @ + temp !
    \ f1 = f2
    f2 @ f1 !
    \ f2 = temp
    temp @ f2 !
  repeat
  sum @ .
;

0 limit !
0 ready !
ei
begin
  ready @
until
di
solve
