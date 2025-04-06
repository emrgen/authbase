import {Box, Flex, Heading, Input, Stack, Field as FormField} from "@chakra-ui/react";
import {Field, Formik} from "formik";
import {authbase} from "../api/client.ts";
import {Button} from "../components/ui/button";
import {useNavigate}  from "react-router";

export function LoginPage() {
  const navigate = useNavigate();
  // const bg = useColorModeValue("white", "gray.700");
  // const toast = useToast();

  // useEffect(() => {
  //   if (user.id) {
  //     navigate("/app");
  //   }
  // }, [ navigate, user ]);

  return (
    <Stack pos={'absolute'} w={'full'} h={'full'} justify={'center'} align={'center'}>
      <Box><Flex
        minH={"100vh"}
        align={"center"}
        justify={"center"}
        // avoid flicker on back press when logged in
        // opacity={user.id ? 0 : 1}
      >
        <Stack gap={8} mx={"auto"} maxW={"lg"} py={12} px={6} w={"450px"}>
          <Stack align={"center"}>
            <Heading fontSize={"4xl"}>Sign in</Heading>
          </Stack>
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
              console.log(project,email, password)
              setSubmitting(true);
              authbase.admin.adminLoginUsingPassword({
                body:{
                  projectName: project,
                  email,
                  password
                }
              }).then((res) => {
                const {data} = res;
                const {token = {}} = data;
                const {accessToken = '', refreshToken = ''} = token;
                localStorage.setItem("accessToken", accessToken.toString());
                localStorage.setItem("refreshToken", refreshToken.toString());
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
            {({
                handleSubmit,
                isSubmitting,
                /* and other goodies */
              }) => {
              return (
                <form onSubmit={handleSubmit}>
                  <Box rounded={"lg"} boxShadow={"lg"} p={8}>
                    <Stack gap={4}>
                      <FormField.Root
                        id="project"
                        // isInvalid={Boolean(errors.email) && touched.email}
                      >
                        <FormField.Label>Project</FormField.Label>
                        <Field
                          as={Input}
                          type="text"
                          id="project"
                          name="project"
                          // isDisabled={isSubmitting}
                          // validate={emailValidator}
                        />
                      </FormField.Root>

                      <FormField.Root
                        id="email"
                        // isInvalid={Boolean(errors.email) && touched.email}
                      >
                        <FormField.Label>Email</FormField.Label>
                        <Field
                          as={Input}
                          type="email"
                          id="email"
                          name="email"
                          // isDisabled={isSubmitting}
                          // validate={emailValidator}
                        />
                      </FormField.Root>

                      <FormField.Root
                        id="password"
                      >
                        <FormField.Label>Password</FormField.Label>
                        <Field
                          as={Input}
                          type="password"
                          id="password"
                          name="password"
                          // isDisabled={isSubmitting}
                          // validate={passwordValidator}
                        />
                      </FormField.Root>
                      <Stack gap={10}>
                        {/* <Stack
                      direction={{ base: "column", sm: "row" }}
                      align={"start"}
                      justify={"space-between"}
                    >
                      <Field
                        as={Checkbox}
                        name="rememberMe"
                        id="rememberMe"
                        isDisabled={isSubmitting}
                      >
                        Remember me
                      </Field>
                    </Stack> */}

                        <Stack gap={4} mt={4}>
                          <Button
                            // _hover={{
                            //   bg: "black",
                            // }}
                            type="submit"
                            // isDisabled={isSubmitting}
                            loading={isSubmitting}
                          >
                            Sign in
                          </Button>
                        </Stack>
                      </Stack>
                    </Stack>
                  </Box>
                </form>
              );
            }}
          </Formik>
        </Stack>
      </Flex>
      </Box>
    </Stack>
  );
}
