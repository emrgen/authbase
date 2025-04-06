import {Box, Flex, Heading, Input, Stack, Text, Button, Field as FormField} from "@chakra-ui/react";
import {Field, Formik} from "formik";
import {useNavigate} from "react-router";

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
              email: "",
              password: "",
            }}
            onSubmit={(
              values: { email: string; password: string },
              {setSubmitting},
            ) => {
              const {email, password} = values;
              console.log(setSubmitting, email, password)
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
                        id="email"
                        // isInvalid={Boolean(errors.email) && touched.email}
                      >
                        <FormField.Label>Email</FormField.Label>
                        <Field
                          as={Input}
                          type="text"
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
                            // isLoading={isSubmitting}
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
