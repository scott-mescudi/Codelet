import { CodeBox } from "@/components/CodeBlock";
import Image from "next/image";

const code = `
query := map[string]string{
  "add_user":             'INSERT INTO users(username, email, role, password_hash) VALUES($1, $2, $3, $4)',
  "get_user_password":    'SELECT password_hash, id FROM users WHERE email=$1',
  "get_user_password_via_id":    'SELECT password_hash FROM users WHERE id=$1',
  "update_user_password": 'UPDATE users SET password_hash=$1 WHERE id=$2',
}

db, err = dataAccess.PrepareStatements(query, db)
if err != nil {
  log.Fatalln(err)
}

srv := userMethods.UserService{Db: db}

app.HandleFunc("/api/v1/register", srv.Signup)
app.HandleFunc("/api/v1/login", srv.Login)
app.Handle("/api/v1/update/password", middleware.AuthMiddleware(srv.ChangePassword)) 
app.Handle("/api/v1/logout", middleware.AuthMiddleware(srv.Logout))  
app.Handle("/api/v1/refresh", middleware.AuthMiddleware(srv.Refresh))

if err := http.ListenAndServe(":8080", app); err != nil {
  log.Fatalln(err)
}
`

const code2 = `
query := map[string]string{
  "add_user":             'INSERT INTO users(username, email, role, password_hash) VALUES($1, $2, $3, $4)',
  "get_user_password":    'SELECT password_hash, id FROM users WHERE email=$1',
  "get_user_password_via_id":    'SELECT password_hash FROM users WHERE id=$1',
  "update_user_password": 'UPDATE users SET password_hash=$1 WHERE id=$2',
}

db, err = dataAccess.PrepareStatements(query, db)


`

export default function Home() {
  return (
    <>
    <div className="sm:w-full flex flex-col gap-5 sm:p-x0 w-full px-5">
      <div className="h-screen w-full bg-neutral-black border border-opacity-15 border-white rounded-2xl"></div>
      <div className="h-screen w-full bg-neutral-black border border-opacity-15 border-white rounded-2xl"></div>
    </div>
    </>
  );
}
