// qp8dp - a pure Go data query parser
//
// Copyright (c) 2022 Michael D Henderson
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package qp8db

//  query     = SELECT selList FROM fromList (WHERE condList SEMICOLON) .
//  selList   = ATTRIBUTE (COMMA ATTRIBUTE)* .
//  fromList  = RELATION (COMMA RELATION)* .
//  condition = (condition AND condition)
//            | (ATTRIBUTE EQUALS ATTRIBUTE)
//            | (ATTRIBUTE IN query)
//            | (ATTRIBUTE LIKE PATTERN) .

type QueryNode struct {
	Select *SelListNode
	From   *FromListNode
	Where  *CondListNode
}

type SelListNode struct {
	Attributes []*AttributeTerm
}

type FromListNode struct {
	Relations []*RelationTerm
}

type CondListNode struct {
	CondList []*ConditionNode
}

type ConditionNode struct {
	And    *AndNode
	Equals *EqualsNode
	In     *InNode
	Like   *LikeNode
}

type AndNode struct {
	LeftSide  *ConditionNode
	RightSide *ConditionNode
}

type EqualsNode struct {
	LeftSide  *AttributeTerm
	RightSide *AttributeTerm
}

type InNode struct {
	Attribute *AttributeTerm
	Query     *QueryNode
}

type LikeNode struct {
	Attribute *AttributeTerm
	Pattern   *PatternTerm
}

type AttributeTerm struct{}

type PatternTerm struct{}

type RelationTerm struct{}
