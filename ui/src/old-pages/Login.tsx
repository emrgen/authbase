import {authbase} from "@/api/client.ts";
import {LoadingButton} from "@/components/ui/loading-button.tsx";
import {Spinner} from "@/components/ui/spinner.tsx";
import {cn} from "@/lib/utils"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {Input} from "@/components/ui/input"
import {Label} from "@/components/ui/label"
import {useUserStore} from "@/store/user.ts";
import {Field, Formik} from "formik";
import {useEffect} from "react";
import {useNavigate} from "react-router";

export function LoginForm() {
  const navigate = useNavigate();
  const setUser = useUserStore((state) => state.setUser);
  const setIsAuthenticated = useUserStore((state) => state.setIsAuthenticated);

  // check if the token is already in local storage and still valid
  useEffect(() => {
    const accessToken = localStorage.getItem("accessToken");
    const refreshToken = localStorage.getItem("refreshToken");

    if (accessToken && refreshToken) {
      authbase.account.getCurrentAccount({}).then((res) => {
        const {data} = res;
        const {account = {}} = data;
        setUser({
          id: account.id!,
          name: account.username!,
          email: account.email!,
        });
        setIsAuthenticated(true);
        navigate("/");
      }).catch((err) => {
        console.log(err);
        localStorage.removeItem("accessToken");
        localStorage.removeItem("refreshToken");
      })
    }
  }, [navigate, setIsAuthenticated, setUser]);

  if (localStorage.getItem("accessToken")) {
    return (
      <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
        <div className="w-full max-w-sm">
          <Spinner/>
        </div>
      </div>
    )
  }

  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <div className={cn("flex flex-col gap-6")}>
          <Card>
            <CardHeader>
              <CardTitle className="text-2xl text-center">Login</CardTitle>
            </CardHeader>
            <CardContent>
              <Formik
                initialValues={{
                  project: "",
                  email: "",
                  password: "",
                }}
                onSubmit={(
                  values: { project: string, email: string; password: string },
                  {setSubmitting},
                ) => {
                  const {project, email, password} = values;
                  console.log(project, email, password)
                  setSubmitting(true);
                  authbase.admin.adminLoginUsingPassword({
                    body: {
                      projectName: project,
                      email,
                      password
                    }
                  }).then((res) => {
                    const {data} = res;
                    const {token = {}, account = {}} = data;
                    const {accessToken = '', refreshToken = ''} = token;
                    localStorage.setItem("accessToken", accessToken.toString());
                    localStorage.setItem("refreshToken", refreshToken.toString());
                    setUser({
                      id: account.id!,
                      name: account.username!,
                      email:account.email!,
                    });
                    setIsAuthenticated(true);

                    navigate("/");
                  }).catch((err) => {
                    console.log(err);
                    // toast({
                    //   title: "Error",
                    //   description: err.response.data.message,
                    //   status: "error",
                    //   duration: 5000,
                    //   isClosable: true,
                    // });
                  }).finally(() => {
                    setSubmitting(false);
                  })
                }}
              >
                {({handleSubmit, isSubmitting}) => {
                  return (
                    <form onSubmit={handleSubmit}>
                      <div className="flex flex-col gap-6">

                        <div className="grid gap-2">
                          <Label htmlFor="project">Project</Label>
                          <Field
                            as={Input}
                            type="text"
                            id="project"
                            name="project"
                            placeholder="Project Name"
                            required
                          />
                        </div>

                        <div className="grid gap-2">
                          <Label htmlFor="email">Email</Label>
                          <Field
                            as={Input}
                            type="email"
                            id="email"
                            name="email"
                            placeholder="m@example.com"
                            required
                          />

                        </div>
                        <div className="grid gap-2">
                          <div className="flex items-center">
                            <Label htmlFor="password">Password</Label>
                            <a
                              href="#"
                              className="ml-auto inline-block text-sm underline-offset-4 hover:underline"
                            >
                              Forgot your password?
                            </a>
                          </div>
                          <Field
                            as={Input}
                            id="password"
                            type="password"
                            name="password"
                            autoComplete="current-password"
                            required
                          />
                        </div>

                        <LoadingButton type="submit" className="w-full" loading={isSubmitting}>
                          Login
                        </LoadingButton>
                        {/*<Button variant="outline" className="w-full">*/}
                        {/*  Login with Google*/}
                        {/*</Button>*/}
                      </div>
                    </form>
                  )
                }
                }
              </Formik>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
