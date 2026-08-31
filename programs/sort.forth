
variable idx
variable ready
variable n0
variable n1
variable n2

: handle_input
  key dup 10 = if
    drop
    1 ready !
  else
    48 -
    dup idx @ 0 = if
      n0 !
    else
      idx @ 1 = if
        n1 !
      else
        n2 !
      then
    then
    idx @ 1+ idx !
  then
  iret
;

: sort3
  n0 @ n1 @ > if
    n0 @ n1 @ n0 ! n1 !
  then
  n1 @ n2 @ > if
    n1 @ n2 @ n1 ! n2 !
  then
  n0 @ n1 @ > if
    n0 @ n1 @ n0 ! n1 !
  then
;

0 idx !
0 ready !
ei
begin
  ready @
until
di
sort3
n0 @ .
n1 @ .
n2 @ .
