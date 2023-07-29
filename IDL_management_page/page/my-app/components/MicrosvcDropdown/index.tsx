import { ContainerVStack, SelectText, AddNewText } from "./styles";
import { Menu, MenuButton, MenuList, MenuItem, Button } from "@chakra-ui/react";
import { ChevronDownIcon } from "@chakra-ui/icons";
import { useRouter } from "next/router";
import { Routes } from "../../constants/Routes.enum";
// import { entryInfoArray } from "../../data/db";
import { Entry } from "../../types/entry";

export default function MicrosvcDropdownComponent(props: any) {
  const router = useRouter();
  const handleCreateSvcClick = (e: any) => {
    e.preventDefault();
    router.push(Routes.CREATE_SERVICE);
  };

  return (
    <>
      <ContainerVStack>
        <SelectText>Select Microservice</SelectText>
        <Menu>
          <MenuButton as={Button} rightIcon={<ChevronDownIcon />}>
            Services
          </MenuButton>
          <MenuList>
            {props.entries.map((entry: Entry, idx: number) => {
              return (
                <MenuItem
                  onClick={() =>
                    router.push(Routes.VIEW_SERVICES + "/" + entry.svcname)
                  }
                  key={idx}
                >
                  {entry.svcname}
                </MenuItem>
              );
            })}
          </MenuList>
        </Menu>

        <p>OR</p>
        <AddNewText onClick={handleCreateSvcClick}>
          Add A New Microservice
        </AddNewText>
      </ContainerVStack>
    </>
  );
}
