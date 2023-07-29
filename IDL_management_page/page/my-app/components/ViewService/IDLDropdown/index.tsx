import { ContainerVStack, SelectText, AddNewText } from "./styles";
import { Menu, MenuButton, MenuList, MenuItem, Button } from "@chakra-ui/react";
import { ChevronDownIcon } from "@chakra-ui/icons";
import { useRouter } from "next/router";
import { Entry } from "../../../types/entry";
import { Routes } from "../../../constants/Routes.enum";
import { Version } from "../../../types/version";

//takes an entry and maps versions to it
export default function IDLDropdownComponent({ id, svcname, versions }: Entry) {
  const router = useRouter();

  return (
    <>
      <ContainerVStack>
        <Menu>
        <MenuButton as={Button} rightIcon={<ChevronDownIcon />}>Select Version</MenuButton>
          <MenuList>
            {versions.map((version: Version, idx: number) => {
              return (
                <MenuItem
                  onClick={() => router.push(Routes.VIEW_SERVICES + "/" + svcname +  "/" + version.vname)}
                  key={idx}
                >
                  {version.vname}
                </MenuItem>
              );
            })}
          </MenuList>
        </Menu>
      </ContainerVStack>
    </>
  );
}
