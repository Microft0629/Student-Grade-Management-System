export namespace model {
	
	export class Course {
	    ID: number;
	    CourseCode: string;
	    CourseName: string;
	    Term: string;
	    Credit: number;
	    Teacher: string;
	    CreatorName: string;
	
	    static createFrom(source: any = {}) {
	        return new Course(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CourseCode = source["CourseCode"];
	        this.CourseName = source["CourseName"];
	        this.Term = source["Term"];
	        this.Credit = source["Credit"];
	        this.Teacher = source["Teacher"];
	        this.CreatorName = source["CreatorName"];
	    }
	}
	export class CourseStatistics {
	    CourseName: string;
	    Term: string;
	    StudentCount: number;
	    AverageScore: number;
	    PassRate: number;
	    HighestScore: number;
	    LowestScore: number;
	    Distribution: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new CourseStatistics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.CourseName = source["CourseName"];
	        this.Term = source["Term"];
	        this.StudentCount = source["StudentCount"];
	        this.AverageScore = source["AverageScore"];
	        this.PassRate = source["PassRate"];
	        this.HighestScore = source["HighestScore"];
	        this.LowestScore = source["LowestScore"];
	        this.Distribution = source["Distribution"];
	    }
	}
	export class ErrorLog {
	    // Go type: time
	    Time: any;
	    Student: string;
	    Course: string;
	    Score: number;
	    Reason: string;
	
	    static createFrom(source: any = {}) {
	        return new ErrorLog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Time = this.convertValues(source["Time"], null);
	        this.Student = source["Student"];
	        this.Course = source["Course"];
	        this.Score = source["Score"];
	        this.Reason = source["Reason"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Student {
	    ID: number;
	    StudentID: string;
	    Name: string;
	    Gender: string;
	    ClassName: string;
	    Major: string;
	
	    static createFrom(source: any = {}) {
	        return new Student(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.StudentID = source["StudentID"];
	        this.Name = source["Name"];
	        this.Gender = source["Gender"];
	        this.ClassName = source["ClassName"];
	        this.Major = source["Major"];
	    }
	}
	export class Grade {
	    ID: number;
	    StudentID: number;
	    CourseID: number;
	    Score: number;
	    GradePoint: number;
	    CreatorName: string;
	    Student: Student;
	    Course: Course;
	
	    static createFrom(source: any = {}) {
	        return new Grade(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.StudentID = source["StudentID"];
	        this.CourseID = source["CourseID"];
	        this.Score = source["Score"];
	        this.GradePoint = source["GradePoint"];
	        this.CreatorName = source["CreatorName"];
	        this.Student = this.convertValues(source["Student"], Student);
	        this.Course = this.convertValues(source["Course"], Course);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class OperationLog {
	    ID: number;
	    // Go type: time
	    Time: any;
	    Operator: string;
	    Action: string;
	    Student: string;
	    StudentID: string;
	    Course: string;
	    Term: string;
	    OldScore: number;
	    NewScore: number;
	    Detail: string;
	
	    static createFrom(source: any = {}) {
	        return new OperationLog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Time = this.convertValues(source["Time"], null);
	        this.Operator = source["Operator"];
	        this.Action = source["Action"];
	        this.Student = source["Student"];
	        this.StudentID = source["StudentID"];
	        this.Course = source["Course"];
	        this.Term = source["Term"];
	        this.OldScore = source["OldScore"];
	        this.NewScore = source["NewScore"];
	        this.Detail = source["Detail"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class StudentPageResult {
	    List: Student[];
	    Total: number;
	    Page: number;
	    PageSize: number;
	
	    static createFrom(source: any = {}) {
	        return new StudentPageResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.List = this.convertValues(source["List"], Student);
	        this.Total = source["Total"];
	        this.Page = source["Page"];
	        this.PageSize = source["PageSize"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StudentRanking {
	    StudentName: string;
	    StudentID: string;
	    TotalScore: number;
	    AverageScore: number;
	    GPA: number;
	    CourseCount: number;
	    Rank: number;
	
	    static createFrom(source: any = {}) {
	        return new StudentRanking(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.StudentName = source["StudentName"];
	        this.StudentID = source["StudentID"];
	        this.TotalScore = source["TotalScore"];
	        this.AverageScore = source["AverageScore"];
	        this.GPA = source["GPA"];
	        this.CourseCount = source["CourseCount"];
	        this.Rank = source["Rank"];
	    }
	}
	export class StudentStatistics {
	    StudentName: string;
	    AverageScore: number;
	    GPA: number;
	    TotalCredits: number;
	    CourseCount: number;
	
	    static createFrom(source: any = {}) {
	        return new StudentStatistics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.StudentName = source["StudentName"];
	        this.AverageScore = source["AverageScore"];
	        this.GPA = source["GPA"];
	        this.TotalCredits = source["TotalCredits"];
	        this.CourseCount = source["CourseCount"];
	    }
	}
	export class User {
	    ID: number;
	    Username: string;
	    Password: string;
	    Role: string;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Username = source["Username"];
	        this.Password = source["Password"];
	        this.Role = source["Role"];
	    }
	}

}

export namespace service {
	
	export class AggregatedGrade {
	    StudentID: string;
	    StudentName: string;
	    CourseName: string;
	    Term: string;
	    Score: number;
	    GradePoint: number;
	    Credit: number;
	
	    static createFrom(source: any = {}) {
	        return new AggregatedGrade(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.StudentID = source["StudentID"];
	        this.StudentName = source["StudentName"];
	        this.CourseName = source["CourseName"];
	        this.Term = source["Term"];
	        this.Score = source["Score"];
	        this.GradePoint = source["GradePoint"];
	        this.Credit = source["Credit"];
	    }
	}
	export class BatchAdjustResult {
	    AffectedCount: number;
	    Details: string[];
	
	    static createFrom(source: any = {}) {
	        return new BatchAdjustResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AffectedCount = source["AffectedCount"];
	        this.Details = source["Details"];
	    }
	}

}

