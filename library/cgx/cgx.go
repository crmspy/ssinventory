package cgx
/*
this is general function collection
it's can use in every time you need
*/

//get offset and limit for paging
func Calcpage(page int)(int,int){
    page -= 1
    limit := 10
    offset := (page * limit) + 1
    return offset,limit
}
