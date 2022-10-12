1.	整体感知
该MatchingService使用Golang编写，提供了http api以提供对已有表信息的匹配查询。该项目的代码设计6个文件，行数大概是600行（含注释，此外另有100行单元测试内容）。在实现过程中，由于本身没有使用过于复杂的算法，所以将重点放在了代码规范性、设计模式和错误检查上。


2.	运行方式
在合适的go环境下（开发时使用的是go 1.18），该demo可以轻松地编译以及运行。其中，服务端信息表的来源是MatchingService/dict.csv；成功运行之后，程序提供一个可以从localhost访问的API，端口号为9527；成功查询之后的返回结果是一个带有列标题的数据集的CSV文件。在使用中，可以通过修改dict.csv以及更换命令进行检验。


3.	实现思路
  API层使用go语言支持的net/http库即可，可以支持路由转发等一系列内容。
  具体的逻辑处理由Matcher进行，Matcher包含的主要成员变量有：
    originalDict: to store the original data in CSV file, for example, [[A,B,C],[a1,b1,c1],[a2,b2,c2]]
    dict: to store a map from title to its column, to help with processing tasks, for example [A:[a1,a2],B:[b1,b2],C:[c1,c2]]
    columnNum: indicates number of columns, since A,B,C, it's 3
    choicesNum: indicates rows in the table(except for title row), it's 2
  Matcher包含的主要成员方法有：
    func (m *Matcher) separateQueries(queries string) ([][]string, error)
    separateQueries: queries like: C1 == "A" or C2 %26= "B" to [[and,C1,==,A][or,C2,&=,B]]
    func (m *Matcher) MatchWithQueries(queries string) ([][]string, error)
    MatchWithQueries: the main function to handle task accepts queries in a whole string and returns answer rows with column name
  其中，separateQueries负责对http请求透传的query字符串做格式化处理，目的是将每一个子请求拆分为以下四元组表答：
    【逻辑连接词（and、or）, column name, operator, value】
  如果不能有效拆分，说明Query有误，需要停止处理，同时返回合适的错误提示信息。
  具体匹配内容的过程由MatchWithQueries函数进行处理。在该过程中，针对每一个四元组query匹配出正确的行，并依次取交集从而最终得到正确结果。

  在对四元组进行处理的过程中，使用了策略模式：Selector接口分别有四种实现，即针对不同的operator进行“分别但又统一”的处理（selectWithQuery 方法）。
  同时，还是在对四元组进行匹配的方式上，使用了小trick，即用map实现了column name到列内容的映射，加快的查找速度，进而提高效率。


4.	规范以及要求
  1.	针对数据源CSV文件，要求：
    a)	至少有一行列标题和一行数据
    b)	标题列的内容仅含有字母和数字
    c)	每一行都是完备的（每一行有相同长度）
  否则，日志打印出具体的错误信息并退出程序。
  2.	针对query的内容。
  有效的http请求格式如：127.0.0.1:9527/?query=C == "c1" or C %26= "c"，其中Query为C == "c1" or C %26= "c"。我们要求：
    a)	Query的每一个元素均由空格分隔
    b)	value的部分必须由双引号包裹起来
    c)	信息需要完整
  否则，停止处理该请求，并返回错误提示信息，如 near XXX。


5.	总结与反思
  1.	该项目的闪光点包括：
    a)	优良的代码习惯：尽管是小项目，但仍有较好的分层、抽取函数意识；
    b)	错误处理，尽可能地进行了错误信息的提示。其中需要指出，对query的错误指出还不够完善，在一些case中不能很好给出修改建议，该改进的逻辑比较复杂，暂时还未做；
    c)	在兼顾代码可读的基础上尽可能使用编程技巧。比如单例模式、策略模式、哈希表。
  2.	未来需要改进的内容：
    a)	更完善的错误提示方式
    b)	改进算法。可以引入倒排索引的概念，记录数据表中每个数据对应的行列坐标，从而提高某些情况下的处理效率。当然，如果是数据相似度比较高的情况下，使用倒排索引的匹配方式可能  反倒不如遍历，针对这一点，可以在初始化的时候引入相似度指标，用以评估使用何种搜索方式。
    c)	计算方式的修改。如果在实际使用场景中，每一次查询中包含的query数量比较多，我们可以采用分布式计算的方式，比如MapReduce。
